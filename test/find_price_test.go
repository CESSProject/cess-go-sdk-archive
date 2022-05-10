package test

import (
	"cess-go-sdk/sdk"
	"testing"
)

func TestFindPrice(t *testing.T) {
	var find sdk.FindSDK
	find.ChainData.CessRpcAddr = ""
	spaceprice, err := find.FindPrice()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(spaceprice)
	}
}
