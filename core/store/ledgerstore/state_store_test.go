package ledgerstore

import (
	"testing"

	"dan-road-vbft/account"
	"dan-road-vbft/common"
	"dan-road-vbft/core/payload"
	"dan-road-vbft/core/states"
	scommon "dan-road-vbft/core/store/common"
	"dan-road-vbft/core/store/statestore"
	"dan-road-vbft/crypto/keypair"
)

func TestContractState(t *testing.T) {
	batch, err := getStateBatch()
	if err != nil {
		t.Errorf("NewStateBatch error %s", err)
		return
	}
	testCode := []byte("testcode")

	deploy := &payload.DeployCode{
		Code:        testCode,
		NeedStorage: false,
		Name:        "testsm",
		Version:     "v1.0",
		Author:      "",
		Email:       "",
		Description: "",
	}

	address := common.AddressFromVmCode(testCode)
	err = batch.TryGetOrAdd(
		scommon.ST_CONTRACT,
		address[:],
		deploy)
	if err != nil {
		t.Errorf("TryGetOrAdd contract error %s", err)
		return
	}

	err = batch.CommitTo()
	if err != nil {
		t.Errorf("batch.CommitTo error %s", err)
		return
	}
	err = testStateStore.CommitTo()
	if err != nil {
		t.Errorf("testStateStore.CommitTo error %s", err)
		return
	}
	contractState1, err := testStateStore.GetContractState(address)
	if err != nil {
		t.Errorf("GetContractState error %s", err)
		return
	}
	if contractState1.Name != deploy.Name ||
		contractState1.Version != deploy.Version ||
		contractState1.Author != deploy.Author ||
		contractState1.Description != deploy.Description ||
		contractState1.Email != deploy.Email {
		t.Errorf("TestContractState failed %+v != %+v", contractState1, deploy)
		return
	}
}

func TestBookkeeperState(t *testing.T) {
	batch, err := getStateBatch()
	if err != nil {
		t.Errorf("NewStateBatch error %s", err)
		return
	}
	acc1 := account.NewAccount("")
	acc2 := account.NewAccount("")

	currBookkeepers := make([]keypair.PublicKey, 0)
	currBookkeepers = append(currBookkeepers, acc1.PublicKey)
	currBookkeepers = append(currBookkeepers, acc2.PublicKey)
	nextBookkeepers := make([]keypair.PublicKey, 0)
	nextBookkeepers = append(nextBookkeepers, acc1.PublicKey)
	nextBookkeepers = append(nextBookkeepers, acc2.PublicKey)

	bookkeeperState := &states.BookkeeperState{
		CurrBookkeeper: currBookkeepers,
		NextBookkeeper: nextBookkeepers,
	}

	batch.TryAdd(scommon.ST_BOOKKEEPER, BOOKKEEPER, bookkeeperState)
	err = batch.CommitTo()
	if err != nil {
		t.Errorf("batch.CommitTo error %s", err)
		return
	}
	err = testStateStore.CommitTo()
	if err != nil {
		t.Errorf("testStateStore.CommitTo error %s", err)
		return
	}
	bookState, err := testStateStore.GetBookkeeperState()
	if err != nil {
		t.Errorf("GetBookkeeperState error %s", err)
		return
	}
	currBookkeepers1 := bookState.CurrBookkeeper
	nextBookkeepers1 := bookState.NextBookkeeper
	for index, pk := range currBookkeepers {
		pk1 := currBookkeepers1[index]
		if !keypair.ComparePublicKey(pk, pk1) {
			t.Errorf("TestBookkeeperState currentBookkeeper failed")
			return
		}
	}
	for index, pk := range nextBookkeepers {
		pk1 := nextBookkeepers1[index]
		if !keypair.ComparePublicKey(pk, pk1) {
			t.Errorf("TestBookkeeperState nextBookkeeper failed")
			return
		}
	}
}

func getStateBatch() (*statestore.StateBatch, error) {
	testStateStore.NewBatch()
	batch := testStateStore.NewStateBatch()
	return batch, nil
}
