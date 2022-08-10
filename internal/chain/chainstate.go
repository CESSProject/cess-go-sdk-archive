package chain

import (
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

//UserHoldSpaceDetails means to get specific information about user space
func (ci *CessInfo) UserHoldSpaceDetails(AccountPublicKey string) (UserHoldSpaceDetails, error) {
	var (
		err  error
		data UserHoldSpaceDetails
	)
	api.getSubstrateApiSafe()
	defer func() {
		api.releaseSubstrateApi()
		err := recover()
		if err != nil {
			fmt.Printf("[Error]Recover UserHoldSpaceDetails panic fail :%s\n", err)
		}
	}()
	meta, err := api.r.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", ci.ChainModule, ci.ChainModuleMethod)
	}

	publickey, err := types.NewMultiAddressFromHexAccountID(AccountPublicKey)
	if err != nil {
		return data, err
	}
	key, err := types.CreateStorageKey(meta, ci.ChainModule, ci.ChainModuleMethod, publickey.AsID[:])
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", ci.ChainModule, ci.ChainModuleMethod)
	}

	_, err = api.r.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", ci.ChainModule, ci.ChainModuleMethod)
	}
	return data, nil
}

func (userinfo UserHoldSpaceDetails) String() string {
	ret := "———————————————————You Purchased Space———————————————————\n"
	ret += "                   PurchasedSpace:" + userinfo.PurchasedSpace.String() + "(KB)\n"
	ret += "                   UsedSpace:" + userinfo.UsedSpace.String() + "(KB)\n"
	ret += "                   RemainingSpace:" + userinfo.RemainingSpace.String() + "(KB)\n"
	ret += "—————————————————————————————————————————————————————————"
	return ret
}

//GetPrice means the size of the space purchased by all customers of the whole CESS system
func (ci *CessInfo) GetPrice() (types.U128, error) {
	var (
		err  error
		data types.U128
	)
	api.getSubstrateApiSafe()
	defer func() {
		api.releaseSubstrateApi()
		err := recover()
		if err != nil {
			fmt.Printf("[Error]Recover UserHoldSpaceDetails panic :%s\n", err)
		}
	}()
	meta, err := api.r.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", ci.ChainModule, ci.ChainModuleMethod)
	}

	key, err := types.CreateStorageKey(meta, ci.ChainModule, ci.ChainModuleMethod)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", ci.ChainModule, ci.ChainModuleMethod)
	}

	_, err = api.r.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", ci.ChainModule, ci.ChainModuleMethod)
	}
	return data, nil
}

//GetFileInfo means to get the specific information of the file through the current fileid
func (ci *CessInfo) GetFileInfo(fileid string) (FileMetaInfo, error) {
	var (
		err  error
		data FileMetaInfo
	)

	api.getSubstrateApiSafe()
	defer func() {
		api.releaseSubstrateApi()
		err := recover()
		if err != nil {
			fmt.Printf("[Error]Recover UserHoldSpaceDetails panic fail :%s\n", err)
		}
	}()
	meta, err := api.r.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", ci.ChainModule, ci.ChainModuleMethod)
	}
	id, err := types.EncodeToBytes(fileid)

	key, err := types.CreateStorageKey(meta, ci.ChainModule, ci.ChainModuleMethod, types.NewBytes(id))
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", ci.ChainModule, ci.ChainModuleMethod)
	}

	_, err = api.r.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.New("Empty")
	}
	return data, nil
}

//GetFileList means to get a list of all files of the current user
func (ci *CessInfo) GetFileList(AccountPublicKey string) ([][]byte, error) {
	var (
		err  error
		data = make([][]byte, 0)
	)
	api.getSubstrateApiSafe()
	defer func() {
		api.releaseSubstrateApi()
		err := recover()
		if err != nil {
			fmt.Printf("[Error]Recover UserHoldSpaceDetails panic fail :%s\n", err)
		}
	}()
	meta, err := api.r.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", ci.ChainModule, ci.ChainModuleMethod)
	}
	publickey, err := types.NewMultiAddressFromHexAccountID(AccountPublicKey)
	if err != nil {
		return data, err
	}

	key, err := types.CreateStorageKey(meta, ci.ChainModule, ci.ChainModuleMethod, publickey.AsID[:])
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", ci.ChainModule, ci.ChainModuleMethod)
	}

	_, err = api.r.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", ci.ChainModule, ci.ChainModuleMethod)
	}
	return data, nil
}

func (fileinfo FileInfo) String() string {
	ret := "———————————————————File Information———————————————————\n"
	ret += fmt.Sprintf("                  Filename:%v\n", string(fileinfo.File_Name[:]))
	ret += fmt.Sprintf("                  Public:%v\n", fileinfo.Public)
	ret += fmt.Sprintf("                  Filehash:%v\n", string(fileinfo.FileHash[:]))
	ret += fmt.Sprintf("                  Backups:%v\n", fileinfo.Backups)
	ret += fmt.Sprintf("                  Filesize:%v\n", fileinfo.FileSize)
	ret += fmt.Sprintf("                  Downloadfee:%v\n", fileinfo.Downloadfee)
	return ret
}

//GetSchedulerInfo means to get all currently registered schedulers
func (ci *CessInfo) GetSchedulerInfo() ([]SchedulerInfo, error) {
	var (
		err  error
		data []SchedulerInfo
	)
	api.getSubstrateApiSafe()
	defer func() {
		api.releaseSubstrateApi()
		err := recover()
		if err != nil {
			fmt.Printf("[Error]Recover UserHoldSpaceDetails panic fail :%s\n", err)
		}
	}()
	meta, err := api.r.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", ci.ChainModule, ci.ChainModuleMethod)
	}

	//publickey, err := types.NewMultiAddressFromHexAccountID(config.ClientConf.ChainData.AccountPublicKey)
	//if err != nil {
	//	return data, err
	//}
	key, err := types.CreateStorageKey(meta, ci.ChainModule, ci.ChainModuleMethod)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", ci.ChainModule, ci.ChainModuleMethod)
	}

	_, err = api.r.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", ci.ChainModule, ci.ChainModuleMethod)
	}
	return data, nil
}

func GetPublicKey(privatekey string) ([]byte, error) {
	kring, err := signature.KeyringPairFromSecret(privatekey, 0)
	if err != nil {
		return nil, err
	}
	return kring.PublicKey, nil
}
