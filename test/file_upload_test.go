package test

import (
	"cess-go-sdk/sdk"
	"testing"
)

func TestFileUpload(t *testing.T) {
	var file sdk.FileSDK
	file.ChainData.CessRpcAddr = ""
	file.ChainData.IdAccountPhraseOrSeed = ""
	file.ChainData.WalletAddress = ""
	//When sending a file, send it as a file block of 2kb
	blocksize := sdk.MB_1 * sdk.BlockSize(1)
	path := ""
	privatekey := ""
	fileid, err := file.FileUpload(blocksize, path, privatekey)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(fileid)
	}
}
