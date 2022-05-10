package test

import (
	"cess-go-sdk/sdk"
	"testing"
)

func TestFindFile(t *testing.T) {
	var find sdk.FindSDK
	find.ChainData.CessRpcAddr = ""
	fileid := ""
	fileinfo, err := find.FindFile(fileid)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(fileinfo)
	}
}

func TestFindFileList(t *testing.T) {
	var find sdk.FindSDK
	find.ChainData.CessRpcAddr = ""
	find.ChainData.AccountPublicKey = ""
	filelist, err := find.FindFileList()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(filelist)
	}
}
