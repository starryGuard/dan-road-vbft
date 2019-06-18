package rest

import (
	"bytes"
	"dan-road-vbft/common"
	"fmt"
	goSdk "github.com/ontio/ontology-go-sdk"
	"testing"
)

func TestSignTX(t *testing.T) {
	sdk := goSdk.NewOntologySdk()
	//sdk.NewRpcClient().SetAddress("http://172.17.0.2:20336")
	wallet, err := sdk.OpenWallet("./wallet.dat")

	from, err := wallet.GetAccountByAddress("ANRryVJESVNodWvtVcrkLqMimzqGe7jBSZ", []byte("1"))
	to, err := wallet.GetAccountByAddress("AGna5UaixJTZcUkijCnM8n8hifmjCySjTc", []byte("2"))

	tx, err := sdk.Native.Ont.NewTransferTransaction(500, 30000, from.Address, to.Address, 1)
	err = sdk.SignToTransaction(tx, from)
	txbf := new(bytes.Buffer)

	txt, err := tx.IntoImmutable()
	txt.Serialize(txbf)
	if err != nil {
	}
	hexCode := common.ToHexString(txbf.Bytes())
	fmt.Println(hexCode)
}
