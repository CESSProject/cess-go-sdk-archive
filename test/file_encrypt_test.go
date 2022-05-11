package test

import (
	"cess-go-sdk/sdk"
	"testing"
)

func TestFileEncrypt(t *testing.T) {
	var file sdk.FileSDK
	encryptpath, savepath, password := "", "", ""
	err := file.FileEncrypt(encryptpath, savepath, password)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success")
	}
}
