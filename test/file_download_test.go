package test

import (
	"cess-go-sdk/sdk"
	"testing"
)

func TestFileDownload(t *testing.T) {
	var file sdk.FileSDK
	file.ChainData.CessRpcAddr = ""
	fileid := ""
	installpath := ""
	err := file.FileDownload(fileid, installpath)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success")
	}
}
