package test

import (
	"cess-go-sdk/sdk"
	"testing"
)

func TestFindPurchasedSpace(t *testing.T) {
	var find sdk.QuerySDK
	find.ChainData.CessRpcAddr = ""
	PurchasedSpace, err := find.QueryPurchasedSpace()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(PurchasedSpace)
	}
}
