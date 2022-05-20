package sdk

import (
	"cess-go-sdk/config"
	"cess-go-sdk/internal/chain"
	"cess-go-sdk/module/result"
	"cess-go-sdk/tools"
	"fmt"
	"github.com/pkg/errors"
	"strconv"
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
	pubkey, err := tools.DecodeToPub(fs.ChainData.WalletAddress, tools.ChainCessTestPrefix)
	if err != nil {
		fmt.Printf("[Error]The wallet address you entered is incorrect, please re-enter:%v\n", err.Error())
		return userinfo, err
	}
	details, err := ci.UserHoldSpaceDetails(tools.PubBytesTo0XString(pubkey))
	if err != nil {
		return userinfo, errors.Wrap(err, "[Error]Get user data fail")
	}
	if details.UsedSpace.Int64()/1024/1024 == 0 && details.UsedSpace.Int64() != 0 {
		details.UsedSpace.SetInt64(1)
	} else {
		details.UsedSpace.SetInt64(details.UsedSpace.Int64() / 1024 / 1024)
	}
	userinfo.PurchasedSpace = strconv.FormatInt(details.PurchasedSpace.Int64()/1024/1024, 10)
	userinfo.UsedSpace = strconv.FormatInt(details.UsedSpace.Int64(), 10)
	userinfo.RemainingSpace = strconv.FormatInt(details.RemainingSpace.Int64()/1024/1024, 10)
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

	ci.ChainModuleMethod = chain.FindPriceModuleMethod
	Price, err := ci.GetPrice()
	if err != nil {
		return 0, errors.Wrap(err, "[Error]Get price fail")
	}
	PerGB, _ := strconv.ParseFloat(fmt.Sprintf("%.12f", float64(Price.Int64()*int64(1024))/float64(1000000000000)), 64)
	if err != nil {
		return 0, errors.Wrap(err, "[Error]Get price fail,wrong chain data")
	}
	return PerGB, nil
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
	pubkey, err := tools.DecodeToPub(fs.ChainData.WalletAddress, tools.ChainCessTestPrefix)
	if err != nil {
		fmt.Printf("[Error]The wallet address you entered is incorrect, please re-enter:%v\n", err.Error())
		return filelist, err
	}

	data, err := ci.GetFileList(tools.PubBytesTo0XString(pubkey))
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
