package test

import (
	"cess-go-sdk/sdk"
	"testing"
)

func TestExpansion(t *testing.T) {
	var trade sdk.PurchaseSDK
	trade.ChainData.CessRpcAddr = ""
	trade.ChainData.IdAccountPhraseOrSeed = ""
	QuantityOfSpaceYouWantToBuy := 1
	MonthsWantToBuy := 1
	ExpectedPrice := 0
	err := trade.Expansion(QuantityOfSpaceYouWantToBuy, MonthsWantToBuy, ExpectedPrice)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("success")
	}
}
