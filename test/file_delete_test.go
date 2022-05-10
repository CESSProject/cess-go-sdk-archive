package test

import (
	"cess-go-sdk/sdk"
	"testing"
)

func TestFileDelete(t *testing.T) {
	var file sdk.FileSDK
	file.ChainData.CessRpcAddr = ""
	file.ChainData.IdAccountPhraseOrSeed = ""
	fileid := ""
	err := file.FileDelete(fileid)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success")
	}
}
