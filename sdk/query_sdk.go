package sdk

import (
	"cess-go-sdk/config"
	"cess-go-sdk/internal/chain"
	"cess-go-sdk/module/result"
	"cess-go-sdk/tools"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type QuerySDK struct {
	config.CessConf
}

type QueryOperate interface {
	QueryPurchasedSpace() (result.UserSpaceDetails, error)
	QueryPrice() (float64, error)
	QueryFile(string) (result.FileInfo, error)
	QueryFileList() ([]result.FindFileList, error)
}

/*
FindPurchasedSpace means to query the space that the current user has purchased and the space that has been used
*/
func (fs QuerySDK) QueryPurchasedSpace() (result.UserSpaceDetails, error) {
	var userinfo result.UserSpaceDetails
	err := chain.Chain_Init(fs.ChainData.CessRpcAddr)
	if err != nil {
		return userinfo, err
	}

	var ci chain.CessInfo
	ci.RpcAddr = fs.ChainData.CessRpcAddr
	ci.ChainModule = chain.PurchasedSpaceChainModule
	ci.ChainModuleMethod = chain.PurchasedPackage
	pubkey, err := tools.DecodeToPub(fs.ChainData.WalletAddress, tools.ChainCessTestPrefix)
	if err != nil {
		errors.Wrap(err, "[Error]The wallet address you entered is incorrect, please re-enter")
		return userinfo, err
	}
	details, err := ci.UserSpacePackage(tools.PubBytesTo0XString(pubkey))
	if err != nil {
		if err.Error() == "Empty" {
			return userinfo, errors.Wrap(err, "Not Found")
		}
		return userinfo, errors.Wrap(err, "[Error]Get user data fail")
	}

	if details.Space.Uint64 != nil {
		userinfo.TotalSpace = fmt.Sprintf("%d", details.Space.Uint64())
		userinfo.UsedSpace = fmt.Sprintf("%d", details.Used_space.Uint64())
		userinfo.RemainingSpace = fmt.Sprintf("%d", details.Remaining_space.Uint64())
	} else {
		userinfo.TotalSpace = details.Space.String()
		userinfo.UsedSpace = details.Used_space.String()
		userinfo.RemainingSpace = details.Remaining_space.String()
	}
	userinfo.Package_type = uint8(details.Package_type)
	userinfo.Start = uint32(details.Start)
	userinfo.Deadline = uint32(details.Deadline)
	userinfo.State = string(details.State)
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
	ci.PublicKeyOfIdentify, _ = chain.GetPublicKey(fs.ChainData.IdAccountPhraseOrSeed)
	ci.ChainModuleMethod = chain.FindFileModuleMethod[0]
	data, err := ci.GetFileInfo(fileid)
	if err != nil {
		return fileinfo, errors.Wrap(err, "[Error]Get file:"+fileid+" info fail")
	}
	if len(data.Names) == 0 {
		err = errors.New("[Fail]This file may have been deleted by someone")
		return fileinfo, err
	}
	for i := 0; i < len(data.Names); i++ {
		if string(ci.PublicKeyOfIdentify) == string(data.Users[i][:]) {
			fileinfo.FileName = string(data.Names[i])
		}
	}

	fileinfo.FileSize = int64(data.FileSize)
	fileinfo.FileState = string(data.FileState)

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
		errors.Wrap(err, "[Error]The wallet address you entered is incorrect, please re-enter")
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
