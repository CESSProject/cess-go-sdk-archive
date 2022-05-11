package test

import (
	"cess-go-sdk/sdk"
	"testing"
)

func TestFileEncrypt(t *testing.T) {
	var file sdk.FileSDK
	decryptpath, savepath, password := "", "", ""
	err := file.FileEncrypt(decryptpath, savepath, password)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success")
	}
}
