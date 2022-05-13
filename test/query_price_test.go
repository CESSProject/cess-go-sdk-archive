package test

import (
	"cess-go-sdk/sdk"
	"testing"
)

func TestFindPrice(t *testing.T) {
	var find sdk.QuerySDK
	find.ChainData.CessRpcAddr = ""
	spaceprice, err := find.QueryPrice()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(spaceprice)
	}
}
