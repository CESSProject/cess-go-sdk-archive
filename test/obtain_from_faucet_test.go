package test

import (
	"cess-go-sdk/sdk"
	"testing"
)

func TestObtainFromFaucet(t *testing.T) {
	var trade sdk.TradeSDK
	trade.CessConf.ChainData.CessRpcAddr = ""
	trade.CessConf.ChainData.FaucetAddress = ""
	pbk := ""
	err := trade.ObtainFromFaucet(pbk)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success")
	}
}
