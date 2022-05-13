package sdk

import (
	"cess-go-sdk/config"
	"cess-go-sdk/internal/chain"
	"cess-go-sdk/module/result"
	"github.com/pkg/errors"
)

type QuerySDK struct {
	config.CessConf
}

type QueryOperate interface {
	QueryPurchasedSpace() (result.UserHoldSpaceDetails, error)
	QueryPrice() (float64, error)
	QueryFile(string) (result.FileInfo, error)
	QueryFileList() ([]result.FindFileList, error)
}

/*
FindPurchasedSpace means to query the space that the current user has purchased and the space that has been used
*/
func (fs QuerySDK) QueryPurchasedSpace() (result.UserHoldSpaceDetails, error) {
	var userinfo result.UserHoldSpaceDetails
	err := chain.Chain_Init(fs.ChainData.CessRpcAddr)
	if err != nil {
		return userinfo, err
	}

	var ci chain.CessInfo
	ci.RpcAddr = fs.ChainData.CessRpcAddr
	ci.ChainModule = chain.PurchasedSpaceChainModule
	ci.ChainModuleMethod = chain.PurchasedSpaceModuleMethod

	details, err := ci.UserHoldSpaceDetails(fs.ChainData.AccountPublicKey)
	if err != nil {
		return userinfo, errors.Wrap(err, "[Error]Get user data fail")
	}
	userinfo.PurchasedSpace = details.PurchasedSpace.String()
	userinfo.UsedSpace = details.UsedSpace.String()
	userinfo.RemainingSpace = details.RemainingSpace.String()
	return userinfo, nil
}

/*
QueryPrice means to get real-time price of storage space
*/
func (fs QuerySDK) QueryPrice() (float64, error) {
	chain.Chain_Init(fs.ChainData.CessRpcAddr)

	var ci chain.CessInfo
	ci.RpcAddr = fs.ChainData.CessRpcAddr
	ci.ChainModule = chain.FindPriceChainModule

	ci.ChainModuleMethod = chain.FindPriceModuleMethod[0]
	AllPurchased, err := ci.GetPurchasedSpace()
	if err != nil {
		return 0, errors.Wrap(err, "[Error]Get all purchased fail")
	}

	ci.ChainModuleMethod = chain.FindPriceModuleMethod[1]
	AllAvailable, err := ci.GetAvailableSpace()
	if err != nil {
		return 0, errors.Wrap(err, "[Error]Get all available fail")
	}

	var purc int64
	var ava int64
	if AllPurchased.Int != nil {
		purc = AllPurchased.Int64()
	}
	if AllAvailable.Int != nil {
		ava = AllAvailable.Int64()
	}
	if purc == ava {
		err = errors.New("[Success]All space has been bought,The current storage price is:+∞ per (MB)")
		return 0, err
	}

	price := (1024 / float64((ava - purc))) * 1000

	return price, nil
}

/*
QueryFile means to query the files uploaded by the current user
fileid:fileid of the file to look for
*/
func (fs QuerySDK) QueryFile(fileid string) (result.FileInfo, error) {
	var fileinfo result.FileInfo
	chain.Chain_Init(fs.ChainData.CessRpcAddr)

	var ci chain.CessInfo
	ci.RpcAddr = fs.ChainData.CessRpcAddr
	ci.ChainModule = chain.FindFileChainModule

	ci.ChainModuleMethod = chain.FindFileModuleMethod[0]
	data, err := ci.GetFileInfo(fileid)
	if err != nil {
		return fileinfo, errors.Wrap(err, "[Error]Get file:"+fileid+" info fail")
	}
	if len(data.File_Name) == 0 {
		err = errors.New("[Fail]This file may have been deleted by someone")
		return fileinfo, err
	}
	fileinfo.FileName = string(data.File_Name[:])
	fileinfo.FileHash = string(data.FileHash[:])
	fileinfo.Public = bool(data.Public)
	fileinfo.Backups = int8(data.Backups)
	fileinfo.FileSize = int64(data.FileSize)
	fileinfo.Downloadfee = data.Downloadfee.Int64()

	return fileinfo, nil
}

func (fs QuerySDK) QueryFileList() ([]result.FindFileList, error) {
	chain.Chain_Init(fs.ChainData.CessRpcAddr)

	var ci chain.CessInfo
	filelist := make([]result.FindFileList, 0)
	ci.RpcAddr = fs.ChainData.CessRpcAddr
	ci.ChainModule = chain.FindFileChainModule
	ci.ChainModuleMethod = chain.FindFileModuleMethod[1]
	data, err := ci.GetFileList(fs.ChainData.AccountPublicKey)
	if err != nil {
		return nil, errors.Wrap(err, "[Error]Get file list fail")
	}
	for _, v := range data {
		var fileresult result.FindFileList
		fileresult.FileId = string(v)
		filelist = append(filelist, fileresult)
	}
	return filelist, nil
}
