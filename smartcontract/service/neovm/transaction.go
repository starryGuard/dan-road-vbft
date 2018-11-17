package neovm

import (
	"dan-road-vbft/core/types"
	vm "dan-road-vbft/vm/neovm"
	vmtypes "dan-road-vbft/vm/neovm/types"
)

// GetExecutingAddress push transaction's hash to vm stack
func TransactionGetHash(service *NeoVmService, engine *vm.ExecutionEngine) error {
	txn, _ := vm.PopInteropInterface(engine)
	tx := txn.(*types.Transaction)
	txHash := tx.Hash()
	vm.PushData(engine, txHash.ToArray())
	return nil
}

// TransactionGetType push transaction's type to vm stack
func TransactionGetType(service *NeoVmService, engine *vm.ExecutionEngine) error {
	txn, _ := vm.PopInteropInterface(engine)
	tx := txn.(*types.Transaction)
	vm.PushData(engine, int(tx.TxType))
	return nil
}

// TransactionGetAttributes push transaction's attributes to vm stack
func TransactionGetAttributes(service *NeoVmService, engine *vm.ExecutionEngine) error {
	vm.PopInteropInterface(engine)
	attributList := make([]vmtypes.StackItems, 0)
	vm.PushData(engine, attributList)
	return nil
}
