package test

import (
	"cess-go-sdk/sdk"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"os"
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
	backups := ""
	privatekey := ""
	fileid, err := file.FileUpload(blocksize, path, backups, privatekey)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(fileid)
	}
}

func TestFileSplit(t *testing.T) {
	const chunkSize = sdk.MB_1 * sdk.BlockSize(1)
	fileInfo, err := os.Stat("D://爱情转移.mp3")
	if err != nil {
		panic(err)
	}

	num := math.Ceil(float64(fileInfo.Size() / chunkSize))

	fi, err := os.OpenFile("D://爱情转移.mp3", os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	b := make([]byte, chunkSize)
	var i int64 = 1
	filehash := ""
	for ; i <= int64(num); i++ {
		fi.Seek((i-1)*chunkSize, 0)
		if len(b) > int(fileInfo.Size()-(i-1)*chunkSize) {
			b = make([]byte, fileInfo.Size()-(i-1)*chunkSize)
		}
		fi.Read(b)
		h := sha256.New()
		h.Write(b)
		filehash += hex.EncodeToString(h.Sum(nil))
	}
	h := sha256.New()
	h.Write([]byte(filehash))
	filehash = "cess" + hex.EncodeToString(h.Sum(nil))
	fi.Close()
}
