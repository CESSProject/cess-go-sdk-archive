package sdk

import (
	"cess-go-sdk/config"
	"cess-go-sdk/internal/chain"
	"cess-go-sdk/internal/fileHandling"
	"cess-go-sdk/internal/rpc"
	"cess-go-sdk/module"
	"cess-go-sdk/tools"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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

	// file meta info
	fmeta, err := ci.GetFileInfo(fileid)
	if err != nil {
		return err
	}

	if string(fmeta.FileState) != "active" {
		return errors.New("[Tips]The file " + fileid + " has not been backed up, please try again later")
	}

	_, err = os.Stat(installpath)
	if err != nil {
		err = os.MkdirAll(installpath, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "[Error]Create install path error")
		}
	}

	for i := 0; i < len(fmeta.ChunkInfo); i++ {
		// Download the file from the scheduler service
		fname := filepath.Join(installpath, string(fmeta.ChunkInfo[i].ChunkId))
		downloadFromStorage(fname, string(fmeta.ChunkInfo[i].MinerIp))
	}

	r := len(fmeta.ChunkInfo) / 3
	d := len(fmeta.ChunkInfo) - r
	err = fileHandling.ReedSolomon_Restore(installpath, fileid, d, r)
	if err != nil {
		return errors.New("[Tips]The file " + fileid + " download failed, please try again later")
	}
	os.Rename(filepath.Join(installpath, fileid), filepath.Join(installpath, string(fmeta.Names[0])))
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

// Download files from cess storage service
func downloadFromStorage(fpath string, mip string) error {
	file, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	var client *rpc.Client

	wsURL := "ws://" + string(base58.Decode(mip))

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	client, err = rpc.DialWebsocket(ctx, wsURL, "")
	if err != nil {
		return err
	}

	var wantfile rpc.FileDownloadReq
	fname := filepath.Base(fpath)

	wantfile.FileId = fmt.Sprintf("%v", fname)
	wantfile.BlockIndex = 1

	reqmsg := rpc.ReqMsg{}
	reqmsg.Method = module.MinerServiceName
	reqmsg.Service = module.DownloadService
	for {
		data, err := proto.Marshal(&wantfile)
		if err != nil {
			return err
		}
		reqmsg.Body = data
		ctx, _ := context.WithTimeout(context.Background(), 90*time.Second)
		resp, err := client.Call(ctx, &reqmsg)
		if err != nil {
			return err
		}

		var respbody rpc.RespBody
		err = proto.Unmarshal(resp.Body, &respbody)
		if err != nil || respbody.Code != 200 {
			return errors.Wrap(err, "[Error]Download file from CESS reply message"+respbody.Msg+",error")
		}
		var blockData rpc.FileDownloadInfo
		err = proto.Unmarshal(respbody.Data, &blockData)
		if err != nil {
			return errors.Wrap(err, "[Error]Download file from CESS error")
		}

		_, err = file.Write(blockData.Data)
		if err != nil {
			return err
		}

		if blockData.BlockIndex == blockData.BlockTotal {
			break
		}
		wantfile.BlockIndex++
	}
	return nil
}
