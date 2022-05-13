package test

import (
	"cess-go-sdk/sdk"
	"testing"
)

func TestFindFile(t *testing.T) {
	var find sdk.QuerySDK
	find.ChainData.CessRpcAddr = ""
	fileid := ""
	fileinfo, err := find.QueryFile(fileid)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(fileinfo)
	}
}

func TestFindFileList(t *testing.T) {
	var find sdk.QuerySDK
	find.ChainData.CessRpcAddr = ""
	find.ChainData.AccountPublicKey = ""
	filelist, err := find.QueryFileList()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(filelist)
	}
}
