package sdk

import (
	"cess-go-sdk/config"
	"cess-go-sdk/internal/chain"
	"cess-go-sdk/internal/rpc"
	"cess-go-sdk/module"
	"cess-go-sdk/tools"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	cesskeyring "github.com/CESSProject/go-keyring"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"storj.io/common/base58"
)

type FileSDK struct {
	config.CessConf
}
type FileOperate interface {
	FileUpload(BlockSize, string, string, string) (string, error)
	FileDownload(string, string) error
	FileDelete(string) error
	FileDecrypt(string, string, string) error
	FileEncrypt(string, string, string) error
	UploadDeclaration(string, string, string) (string, error)
}

type BlockSize = int64

const (
	KB_1 = BlockSize(1024)
	MB_1 = 1024 * KB_1
	GB_1 = 1024 * MB_1
)

/*
FileUpload means upload files to CESS system
path:The absolute path of the file to be uploaded
backups:Number of backups of files that need to be uploaded
privatekey:Encrypted password for uploaded files
*/
func (fs FileSDK) FileUpload(block BlockSize, fpath, privatekey string) (string, error) {
	blocksize := int(block)
	err := chain.Chain_Init(fs.ChainData.CessRpcAddr)
	if err != nil {
		return "", err
	}
	fstat, err := os.Stat(fpath)
	if err != nil {
		return "", errors.Wrap(err, "[Error]Please enter the correct file path")
	}

	if fstat.IsDir() {
		return "", errors.Wrap(err, "[Error]Please do not upload the folder")
	}

	//Calc file hash
	hash, err := tools.CalcFileHashByChunks(fpath, 1024*1024*1024)
	if err != nil {
		return "", errors.Wrap(err, "[Error]There is a problem with the file, please replace it")
	}
	fileid := "cess" + hash

	if len(privatekey) != 16 && len(privatekey) != 24 && len(privatekey) != 32 && len(privatekey) != 0 {
		return fileid, errors.New("[Error]The password must be 16,24,32 bits long")
	}

	var ci chain.CessInfo
	ci.RpcAddr = fs.ChainData.CessRpcAddr
	ci.ChainModule = chain.FindSchedulerInfoModule
	ci.ChainModuleMethod = chain.FindSchedulerInfoMethod
	ci.TransactionName = chain.UploadDeclaration
	ci.IdentifyAccountPhrase = fs.ChainData.IdAccountPhraseOrSeed
	ci.PublicKeyOfIdentify, err = chain.GetPublicKey(fs.ChainData.IdAccountPhraseOrSeed)
	if err != nil {
		return "", errors.Wrap(err, "[Error]private key error")
	}
	txhash, err := ci.UploadDeclaration(fileid, fstat.Name())
	if txhash == "" {
		return "", errors.Wrap(err, "[Error]UploadDeclaration error")
	}

	var encodefile []byte
	var fsize = int(fstat.Size())
	f, err := os.Open(fpath)
	if err != nil {
		return "", errors.Wrap(err, "[Error]This file was broken")
	}
	filebyte, err := ioutil.ReadAll(f)
	if err != nil {
		f.Close()
		return "", errors.Wrap(err, "[Error]analyze this file error")
	}
	f.Close()
	if len(privatekey) != 0 {
		encodefile, err = tools.AesEncrypt(filebyte, []byte(privatekey))
		if err != nil {
			return fileid, errors.Wrap(err, "[Error]Encode the file fail ,error")
		}
		fsize = len(encodefile)
	}

	var authreq rpc.AuthReq
	authreq.FileId = fileid
	authreq.FileName = fstat.Name()
	authreq.FileSize = uint64(fsize)
	authreq.BlockTotal = uint32(fsize / blocksize)
	if fsize%blocksize != 0 {
		authreq.BlockTotal += 1
	}
	authreq.PublicKey = ci.PublicKeyOfIdentify

	authreq.Msg = []byte(tools.GetRandomcode(16))
	kr, _ := cesskeyring.FromURI(fs.ChainData.IdAccountPhraseOrSeed, cesskeyring.NetSubstrate{})
	// sign message
	sign, err := kr.Sign(kr.SigningContext(authreq.Msg))
	if err != nil {
		return "", errors.Wrap(err, "[Error]Calc sign error")
	}
	authreq.Sign = sign[:]

	schds, err := ci.GetSchedulerInfo()
	if err != nil {
		return "", errors.Wrap(err, "[Error]Get scheduler randomly error")
	}

	//auth
	var client *rpc.Client
	for i, schd := range schds {
		wsURL := "ws://" + string(base58.Decode(string(schd.Ip)))
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		client, err = rpc.DialWebsocket(ctx, wsURL, "")
		defer cancel()
		if err != nil {
			err = errors.New("Connect with scheduler timeout")
			if i == len(schds)-1 {
				return fileid, errors.Wrap(err, "[Error]All scheduler is offline")
			}
			continue
		} else {
			break
		}
	}

	bob, err := proto.Marshal(&authreq)
	if err != nil {
		return "", errors.Wrap(err, "[Error]")
	}
	data, code, err := WriteData(client, module.CtlServiceName, module.UploadAuth, bob)
	if err != nil {
		return "", errors.Wrap(err, "[Error]")
	}
	if code == 201 {
		return fileid, nil
	}
	if code != 200 {
		return "", errors.Errorf("[Error]Auth return code %v", code)
	}

	var filereq rpc.FileUploadReq
	var n int
	var buf = make([]byte, blocksize)
	filereq.Auth = data
	if len(privatekey) == 0 {
		f, err = os.Open(fpath)
		if err != nil {
			return "", errors.Wrap(err, "[Error]This file was broken")
		}
	}
	for i := 0; i < int(authreq.BlockTotal); i++ {
		filereq.BlockIndex = uint32(i + 1)
		if len(privatekey) == 0 {
			f.Seek(int64(i*blocksize), 0)
			n, _ = f.Read(buf)
			filereq.FileData = buf[:n]
		} else {
			if (i+1)*blocksize > fsize {
				filereq.FileData = encodefile[i*blocksize:]
			} else {
				filereq.FileData = encodefile[i*blocksize : (i+1)*blocksize]
			}
		}

		bob, err := proto.Marshal(&filereq)
		if err != nil {
			return "", errors.Wrap(err, "[Error]")
		}

		_, _, err = WriteData(client, module.CtlServiceName, module.UploadService, bob)
		if err != nil {
			return "", errors.Wrap(err, "[Error]")
		}
	}
	return fileid, nil
}

/*
FileDownload means download file by file id
fileid:fileid of the file to download
*/
func (fs FileSDK) FileDownload(fileid, installpath string) error {
	err := chain.Chain_Init(fs.ChainData.CessRpcAddr)
	if err != nil {
		return err
	}
	var ci chain.CessInfo
	ci.RpcAddr = fs.ChainData.CessRpcAddr
	ci.ChainModule = chain.FindFileChainModule
	ci.ChainModuleMethod = chain.FindFileModuleMethod[0]
	fileinfo, err := ci.GetFileInfo(fileid)
	if err != nil {
		return errors.Wrap(err, "[Error]Get file: info fail")
	}
	if fileinfo.File_Name == nil {
		err = errors.New("[Tips]The fileid " + fileid + " has been deleted,the file does not exist")
		return err
	}
	if string(fileinfo.FileState) != "active" {
		err = errors.New("[Tips]The file " + fileid + " has not been backed up, please try again later")
		return err
	}

	_, err = os.Stat(installpath)
	if err != nil {
		err = os.Mkdir(installpath, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "[Error]Create install path error")
		}
	}
	_, err = os.Create(filepath.Join(installpath, string(fileinfo.File_Name[:])))
	if err != nil {
		return errors.Wrap(err, "[Error]Create installed file error ")
	}
	installfile, err := os.OpenFile(filepath.Join(installpath, string(fileinfo.File_Name[:])), os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return errors.Wrap(err, "[Error]:Failed to save key error")
	}
	defer installfile.Close()

	ci.RpcAddr = fs.ChainData.CessRpcAddr
	ci.ChainModule = chain.FindSchedulerInfoModule
	ci.ChainModuleMethod = chain.FindSchedulerInfoMethod
	schds, err := ci.GetSchedulerInfo()
	if err != nil {
		return errors.Wrap(err, "[Error]Get scheduler list error")
	}

	var client *rpc.Client
	for i, schd := range schds {
		wsURL := "ws://" + string(base58.Decode(string(schd.Ip)))
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		client, err = rpc.DialWebsocket(ctx, wsURL, "")
		defer cancel()
		if err != nil {
			err = errors.New("Connect with scheduler timeout")
			if i == len(schds)-1 {
				return errors.Wrap(err, "[Error]All scheduler is offline")
			}
			continue
		} else {
			break
		}
	}

	var wantfile module.FileDownloadReq
	sp := sync.Pool{
		New: func() interface{} {
			return &rpc.ReqMsg{}
		},
	}
	wantfile.FileId = fileid
	wantfile.WalletAddress = fs.ChainData.WalletAddress
	wantfile.BlockIndex = 1

	for {
		data, err := proto.Marshal(&wantfile)
		if err != nil {
			return errors.Wrap(err, "[Error]Marshal req file error")
		}
		req := sp.Get().(*rpc.ReqMsg)
		req.Method = module.DownloadService
		req.Service = module.CtlServiceName
		req.Body = data

		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		resp, err := client.Call(ctx, req)
		cancel()
		if err != nil {
			return errors.Wrap(err, "[Error]Download file fail error")
		}

		var respbody rpc.RespBody
		err = proto.Unmarshal(resp.Body, &respbody)
		if err != nil || respbody.Code != 200 {
			if err != nil {
				return errors.Wrap(err, "[Error]Download file from CESS reply message"+respbody.Msg+",error")
			}
			if err == nil {
				return errors.New("[Error]Download file from CESS reply message" + respbody.Msg)
			}
		}
		var blockData module.FileDownloadInfo
		err = proto.Unmarshal(respbody.Data, &blockData)
		if err != nil {
			return errors.Wrap(err, "[Error]Download file from CESS error")
		}

		_, err = installfile.Write(blockData.Data)
		if err != nil {
			return errors.Wrap(err, "[Error]:Failed to write file's block to file error")
		}

		wantfile.BlockIndex++
		sp.Put(req)
		if blockData.BlockIndex == blockData.BlockTotal {
			break
		}
	}

	return nil
}

/*
FileDelete means to delete the file from the CESS system by the file id
fileid:fileid of the file that needs to be deleted
*/
func (fs FileSDK) FileDelete(fileid string) error {
	err := chain.Chain_Init(fs.ChainData.CessRpcAddr)
	if err != nil {
		return err
	}
	var ci chain.CessInfo
	ci.RpcAddr = fs.ChainData.CessRpcAddr
	ci.IdentifyAccountPhrase = fs.ChainData.IdAccountPhraseOrSeed
	ci.TransactionName = chain.DeleteFileTransactionName

	err = ci.DeleteFileOnChain(fileid)
	if err != nil {
		return errors.Wrap(err, "[Error]Delete file error")
	}
	return nil
}

/*
When you download the file if it is not decrypt, you can decrypt it this way
*/
func (fs FileSDK) FileDecrypt(decryptpath, savepath, password string) error {
	if len(password) != 16 && len(password) != 24 && len(password) != 32 {
		return errors.New("[Error]The password must be 16,24,32 bits long")
	}
	_, err := os.Stat(decryptpath)
	if err != nil {
		_ = errors.Wrap(err, "[Error]There is no such file, please confirm the correct location of the file,and enter the absolute path of the file")
		return err
	}

	encodefile, err := ioutil.ReadFile(decryptpath)
	if err != nil {
		return errors.Wrap(err, "[Error]Failed to read file, please check file integrity")
	}

	decryptfile, err := tools.AesDecrypt(encodefile, []byte(password))
	if err != nil {
		return errors.Wrap(err, "[Error]File decrypt failed, please check your password! error")
	}
	filename := filepath.Base(decryptpath)
	//The decrypted file is saved to the download folder, if the name is the same, the original file will be deleted
	if decryptpath == filepath.Join(savepath, filename) {
		err = os.Remove(decryptpath)
		if err != nil {
			return errors.Wrap(err, "[Error]An error occurred while saving the decrypted file! error")
		}
	}
	fileinfo, err := os.Create(filepath.Join(savepath, filename))
	if err != nil {
		return errors.Wrap(err, "[Error]An error occurred while saving the decrypted file! error")
	}
	defer fileinfo.Close()
	_, err = fileinfo.Write(decryptfile)
	if err != nil {
		return errors.Wrap(err, "[Error]Failed to save decrypted content to file! error")
	}

	return nil
}

/*
You can encrypt files yourself
*/
func (fs FileSDK) FileEncrypt(encryptpath, savepath, password string) error {
	if len(password) != 16 && len(password) != 24 && len(password) != 32 {
		return errors.New("[Error]The password must be 16,24,32 bits long")
	}
	_, err := os.Stat(encryptpath)
	if err != nil {
		_ = errors.Wrap(err, "[Error]There is no such file, please confirm the correct location of the file,and enter the absolute path of the file")
		return err
	}
	filebyte, err := ioutil.ReadFile(encryptpath)
	if err != nil {
		return errors.Wrap(err, "[Error]This file was broken")
	}
	encryptfile, err := tools.AesEncrypt(filebyte, []byte(password))
	if err != nil {
		return errors.Wrap(err, "[Error]encrypt the file fail ,error")
	}
	filename := filepath.Base(encryptpath)
	//The encrypted file is saved to the download folder, if the name is the same, the original file will be deleted
	if encryptpath == filepath.Join(savepath, filename) {
		err = os.Remove(encryptpath)
		if err != nil {
			return errors.Wrap(err, "[Error]An error occurred while saving the encrypted file! error")
		}
	}
	fileinfo, err := os.Create(filepath.Join(savepath, filename))
	if err != nil {
		return errors.Wrap(err, "[Error]An error occurred while saving the encrypted file! error")
	}
	defer fileinfo.Close()
	_, err = fileinfo.Write(encryptfile)
	if err != nil {
		return errors.Wrap(err, "[Error]Failed to save encrypted content to file! error")
	}

	return nil
}

func WriteData(cli *rpc.Client, service, method string, body []byte) ([]byte, int32, error) {
	req := &rpc.ReqMsg{
		Service: service,
		Method:  method,
		Body:    body,
	}
	ctx, _ := context.WithTimeout(context.Background(), 90*time.Second)
	resp, err := cli.Call(ctx, req)
	if err != nil {
		return nil, 0, errors.Wrap(err, "Call err:")
	}

	var b rpc.RespBody
	err = proto.Unmarshal(resp.Body, &b)
	if err != nil {
		return nil, 0, errors.Wrap(err, "Unmarshal:")
	}
	return b.Data, b.Code, err
}
