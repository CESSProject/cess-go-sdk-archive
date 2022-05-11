package test

import (
	"cess-go-sdk/sdk"
	"testing"
)

func TestFileDecrypt(t *testing.T) {
	var file sdk.FileSDK
	decryptpath, savepath, password := "", "", ""
	err := file.FileDecrypt(decryptpath, savepath, password)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success")
	}
}
