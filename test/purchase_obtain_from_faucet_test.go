package test

import (
	"cess-go-sdk/sdk"
	"testing"
)

func TestObtainFromFaucet(t *testing.T) {
	var trade sdk.PurchaseSDK
	trade.CessConf.ChainData.FaucetAddress = ""
	WalletAddress := ""
	err := trade.ObtainFromFaucet(WalletAddress)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success")
	}
}
