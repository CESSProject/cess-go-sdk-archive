package test

import (
	"cess-go-sdk/sdk"
	"testing"
)

func TestFindPurchasedSpace(t *testing.T) {
	var find sdk.FindSDK
	find.ChainData.CessRpcAddr = ""
	find.ChainData.AccountPublicKey = ""
	PurchasedSpace, err := find.FindPurchasedSpace()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(PurchasedSpace)
	}
}
