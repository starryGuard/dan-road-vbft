package test

import (
	"dan-road-vbft/core/types"
	"dan-road-vbft/smartcontract"
	"dan-road-vbft/vm/neovm"
	"os"
	"testing"
)

func TestPackCrash(t *testing.T) {
	// define a leaf
	byteCode := []byte{byte(neovm.PUSH0)}

	// build tree using array packing
	for i := 0; i < 10000; i++ {
		byteCode = append(byteCode, []byte{
			byte(neovm.DUP),
			byte(neovm.PUSH2),
			byte(neovm.PACK),
		}...)
	}
	// compare trees
	byteCode = append(byteCode, []byte{
		byte(neovm.DUP),
		byte(neovm.EQUAL),
	}...)
	// setup VM
	dbFile := "test"
	os.RemoveAll(dbFile)
	//testLevelDB, err := leveldbstore.NewLevelDBStore(dbFile)
	//if err != nil {
	//	panic(err)
	//}
	//store := statestore.NewMemDatabase()
	//testBatch := statestore.NewStateStoreBatch(store, testLevelDB)
	config := &smartcontract.Config{
		Time:   10,
		Height: 10,
		Tx:     &types.Transaction{},
	}
	//cache := storage.NewCloneCache(testBatch)
	sc := smartcontract.SmartContract{
		Config:  config,
		Gas:     200,
		CacheDB: nil,
	}
	engine, err := sc.NewExecuteEngine(byteCode)
	if err != nil {
		panic(err)
		// cause the VM to hang forever
		_, err = engine.Invoke()
		if err != nil {
		}
		panic(err)
	}
}
