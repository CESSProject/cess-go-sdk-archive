package test

import (
	"cess-go-sdk/sdk"
	"testing"
)

func TestObtainFromFaucet(t *testing.T) {
	var trade sdk.TradeSDK
	trade.CessConf.ChainData.FaucetAddress = ""
	AccountPublicKey := ""
	err := trade.ObtainFromFaucet(AccountPublicKey)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success")
	}
}
