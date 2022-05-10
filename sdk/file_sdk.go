package sdk

import (
	"cess-go-sdk/config"
	"cess-go-sdk/internal/chain"
	"cess-go-sdk/internal/rpc"
	"cess-go-sdk/module"
	"cess-go-sdk/tools"
	"context"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type FileSDK struct {
	config.CessConf
}
type FileOperate interface {
	FileUpload(blocksize BlockSize, path, backups, privatekey string) (string, error)
	FileDownload(fileid string) error
	FileDelete(fileid string) error
	FileDecode(path string) error
}

type BlockSize = int64

const (
	B_1  = BlockSize(1024)
	KB_1 = 1024 * B_1
	MB_1 = 1024 * KB_1
	GB_1 = 1024 * MB_1
)

/*
FileUpload means upload files to CESS system
path:The absolute path of the file to be uploaded
backups:Number of backups of files that need to be uploaded
privatekey:Encrypted password for uploaded files
*/
func (fs FileSDK) FileUpload(block BlockSize, path, backups, privatekey string) (string, error) {
	blocksize := int(block)
	err := chain.Chain_Init(fs.ChainData.CessRpcAddr)
	if err != nil {
		return "", err
	}
	file, err := os.Stat(path)
	if err != nil {
		return "", errors.Wrap(err, "[Error]Please enter the correct file path")
	}

	if file.IsDir() {
		return "", errors.Wrap(err, "[Error]Please do not upload the folder")
	}

	spares, err := strconv.Atoi(backups)
	if err != nil {
		return "", errors.Wrap(err, "[Error]Please enter a correct integer")

	}

	filehash, err := tools.CalcFileHash(path)
	if err != nil {
		return "", errors.Wrap(err, "[Error]There is a problem with the file, please replace it")
	}

	fileid, err := tools.GetGuid(1)
	if err != nil {
		return "", errors.Wrap(err, "[Error]Create snowflake fail")
	}
	var blockinfo module.FileUploadInfo
	blockinfo.Backups = backups
	blockinfo.FileId = fileid
	blockinfo.BlockSize = int32(file.Size())
	blockinfo.FileHash = filehash

	blocktotal := 0

	f, err := os.Open(path)
	if err != nil {
		return "", errors.Wrap(err, "[Error]This file was broken")
	}
	defer f.Close()
	filebyte, err := ioutil.ReadAll(f)
	if err != nil {
		return "", errors.Wrap(err, "[Error]analyze this file error")
	}

	var ci chain.CessInfo
	ci.RpcAddr = fs.ChainData.CessRpcAddr
	ci.ChainModule = chain.FindSchedulerInfoModule
	ci.ChainModuleMethod = chain.FindSchedulerInfoMethod
	schds, err := ci.GetSchedulerInfo()
	if err != nil {
		return "", errors.Wrap(err, "[Error]Get scheduler randomly error")
	}
	filesize := new(big.Int)
	fee := new(big.Int)

	ci.IdentifyAccountPhrase = fs.ChainData.IdAccountPhraseOrSeed
	ci.TransactionName = chain.UploadFileTransactionName

	if file.Size()/1024 == 0 {
		filesize.SetInt64(1)
	} else {
		filesize.SetInt64(file.Size() / 1024)
	}
	fee.SetInt64(int64(0))

	_, err = ci.UploadFileMetaInformation(fileid, file.Name(), filehash, privatekey == "", uint8(spares), filesize, fee, fs.ChainData.WalletAddress)
	if err != nil {
		return "", errors.Wrap(err, "[Error]Upload file meta information error")
	}

	var client *rpc.Client
	for i, schd := range schds {
		wsURL := "ws://" + string(base58.Decode(string(schd.Ip)))
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		client, err = rpc.DialWebsocket(ctx, wsURL, "")
		defer cancel()
		if err != nil {
			err = errors.New("Connect with scheduler timeout")
			fmt.Printf("%s[Tips]%sdialog with scheduler:%s fail! reason:%s\n", tools.Yellow, tools.Reset, string(base58.Decode(string(schd.Ip))), err)
			if i == len(schds)-1 {
				return fileid, errors.Wrap(err, "[Error]All scheduler is offline")
			}
			continue
		} else {
			break
		}
	}
	sp := sync.Pool{
		New: func() interface{} {
			return &rpc.ReqMsg{}
		},
	}
	commit := func(num int, data []byte) error {
		blockinfo.BlockNum = int32(num) + 1
		blockinfo.Data = data
		info, err := proto.Marshal(&blockinfo)
		if err != nil {
			return errors.Wrap(err, "[Error]Serialization error, please upload again")
		}
		reqmsg := sp.Get().(*rpc.ReqMsg)
		reqmsg.Body = info
		reqmsg.Method = module.UploadService
		reqmsg.Service = module.CtlServiceName

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		resp, err := client.Call(ctx, reqmsg)
		defer cancel()
		if err != nil {
			return errors.Wrap(err, "[Error]Failed to transfer file to scheduler,error")
		}

		var res rpc.RespBody
		err = proto.Unmarshal(resp.Body, &res)
		if err != nil {
			return errors.Wrap(err, "[Error]Error getting reply from schedule, transfer failed")
		}
		if res.Code != 0 {
			err = errors.New(res.Msg)
			return errors.Wrap(err, "[Error]Upload file fail!scheduler problem")
		}
		sp.Put(reqmsg)
		return nil
	}

	if len(privatekey) != 0 {

		encodefile, err := tools.AesEncrypt(filebyte, []byte(privatekey))
		if err != nil {
			return fileid, errors.Wrap(err, "[Error]Encode the file fail ,error")
		}
		blocks := len(encodefile) / blocksize
		if len(encodefile)%blocksize == 0 {
			blocktotal = blocks
		} else {
			blocktotal = blocks + 1
		}
		blockinfo.Blocks = int32(blocktotal)
		var bar tools.Bar
		bar.NewOption(0, int64(blocktotal))
		for i := 0; i < blocktotal; i++ {
			block := make([]byte, 0)
			if blocks != i {
				block = encodefile[i*blocksize : (i+1)*blocksize]
				bar.Play(int64(i + 1))
			} else {
				block = encodefile[i*blocksize:]
				bar.Play(int64(i + 1))
			}
			err = commit(i, block)
			if err != nil {
				bar.Finish()
				return fileid, errors.Wrap(err, "[Error]:Failed to upload the file error")
			}
		}
		bar.Finish()
	} else {
		fmt.Printf("%s[Tips]%s:upload file:%s without private key", tools.Yellow, tools.Reset, path)
		blocks := len(filebyte) / blocksize
		if len(filebyte)%blocksize == 0 {
			blocktotal = blocks
		} else {
			blocktotal = blocks + 1
		}
		blockinfo.Blocks = int32(blocktotal)
		var bar tools.Bar
		bar.NewOption(0, int64(blocktotal))
		for i := 0; i < blocktotal; i++ {
			block := make([]byte, 0)
			if blocks != i {
				block = filebyte[i*blocksize : (i+1)*blocksize]
				bar.Play(int64(i + 1))
			} else {
				block = filebyte[i*blocksize:]
				bar.Play(int64(i + 1))
			}
			err = commit(i, block)
			if err != nil {
				bar.Finish()
				return fileid, errors.Wrap(err, "[Error]:Failed to upload the file error")
			}
		}
		bar.Finish()
	}
	fmt.Printf("%s[Success]%s:upload file:%s successful!", tools.Green, tools.Reset, path)
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
	if string(fileinfo.FileState) != "active" {
		err = errors.New("[Tips]The file " + fileid + " has not been backed up, please try again later")
		return err
	}
	if fileinfo.File_Name == nil {
		err = errors.New("[Tips]The fileid " + fileid + " used to find the file is incorrect, please try again")
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
			fmt.Printf("%s[Tips]%sdialog with scheduler:%s fail! reason:%s\n", tools.Yellow, tools.Reset, string(base58.Decode(string(schd.Ip))), err)
			if i == len(schds)-1 {
				return errors.Wrap(err, "[Error]All scheduler is offline")
			}
			continue
		} else {
			break
		}
	}

	var wantfile module.FileDownloadReq
	var bar tools.Bar
	var getAllBar sync.Once
	sp := sync.Pool{
		New: func() interface{} {
			return &rpc.ReqMsg{}
		},
	}
	wantfile.FileId = fileid
	wantfile.WalletAddress = fs.ChainData.AccountPublicKey
	wantfile.Blocks = 1

	for {
		data, err := proto.Marshal(&wantfile)
		if err != nil {
			return errors.Wrap(err, "[Error]Marshal req file error")
		}
		req := sp.Get().(*rpc.ReqMsg)
		req.Method = module.DownloadService
		req.Service = module.CtlServiceName
		req.Body = data

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		resp, err := client.Call(ctx, req)
		cancel()
		if err != nil {
			return errors.Wrap(err, "[Error]Download file fail error")
		}

		var respbody rpc.RespBody
		err = proto.Unmarshal(resp.Body, &respbody)
		if err != nil || respbody.Code != 0 {
			return errors.Wrap(err, "[Error]Download file from CESS reply message"+respbody.Msg+",error")
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

		getAllBar.Do(func() {
			bar.NewOption(0, int64(blockData.BlockNum))
		})
		bar.Play(int64(blockData.Blocks))
		wantfile.Blocks++
		sp.Put(req)
		if blockData.Blocks == blockData.BlockNum {
			break
		}
	}

	bar.Finish()
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
When you download the file if it is not decode, you can decode it this way
*/
func (fs FileSDK) FileDecode(decodepath, savepath, password string) error {
	_, err := os.Stat(decodepath)
	if err != nil {
		_ = errors.Wrap(err, "[Error]There is no such file, please confirm the correct location of the file, please enter the absolute path of the file")
		return err
	}

	//fmt.Println("Please enter the file's password:")
	//fmt.Print(">")
	//psw, _ := gopass.GetPasswdMasked()
	encodefile, err := ioutil.ReadFile(decodepath)
	if err != nil {
		return errors.Wrap(err, "[Error]Failed to read file, please check file integrity")
	}

	decodefile, err := tools.AesDecrypt(encodefile, []byte(password))
	if err != nil {
		return errors.Wrap(err, "[Error]File decode failed, please check your password! error")
	}
	filename := filepath.Base(decodepath)
	//The decoded file is saved to the download folder, if the name is the same, the original file will be deleted
	if decodepath == filepath.Join(savepath, filename) {
		err = os.Remove(decodepath)
		if err != nil {
			return errors.Wrap(err, "[Error]An error occurred while saving the decoded file! error")
		}
	}
	fileinfo, err := os.Create(filepath.Join(savepath, filename))
	if err != nil {
		return errors.Wrap(err, "[Error]An error occurred while saving the decoded file! error")
	}
	defer fileinfo.Close()
	_, err = fileinfo.Write(decodefile)
	if err != nil {
		return errors.Wrap(err, "[Error]Failed to save decrypted content to file! error")
	}

	return nil
}
