package account

import (
	"encoding/hex"
	"testing"
)

var id = "did:ont:TSS6S4Xhzt5wtvRBTm4y3QCTRqB4BnU7vT"

func TestCreate(t *testing.T) {
	nonce, _ := hex.DecodeString("4c6b58adc6b8c6774eee0eb07dac4e198df87aae28f8932db3982edf3ff026e4")
	id1, err := CreateID(nonce)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("result ID:", id1)
	if id != id1 {
		t.Fatal("expected ID:", id)
	}
}

func TestVerify(t *testing.T) {
	t.Log("verify", id)
	if !VerifyID(id) {
		t.Error("error: failed")
	}

	invalid := []string{
		"did:ont:",
		"did:else:TSS6S4Xhzt5wtvRBTm4y3QCTRqB4BnU7vT",
		"TSS6S4Xhzt5wtvRBTm4y3QCTRqB4BnU7vT",
		"did:else:TSS6S4Xhzt5wtvRBTm4y3QCT",
		"did:ont:TSS6S4Xhzt5wtvRBTm4y3QCTRqB4BnU7vt",
	}

	for _, v := range invalid {
		t.Log("verify", v)
		if VerifyID(v) {
			t.Error("error: passed")
		}
	}
}
