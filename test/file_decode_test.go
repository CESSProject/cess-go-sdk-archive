package test

import (
	"cess-go-sdk/sdk"
	"testing"
)

func TestFileDecode(t *testing.T) {
	var file sdk.FileSDK
	decodepath, savepath, password := "", "", ""
	err := file.FileDecode(decodepath, savepath, password)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success")
	}
}
