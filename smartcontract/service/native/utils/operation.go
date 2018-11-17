package utils

import (
	"dan-road-vbft/common"
	"dan-road-vbft/errors"
	"dan-road-vbft/smartcontract/event"
	"dan-road-vbft/smartcontract/service/native"
)

func AddCommonEvent(native *native.NativeService, contract common.Address, name string, params interface{}) {
	native.Notifications = append(native.Notifications,
		&event.NotifyEventInfo{
			ContractAddress: contract,
			States:          []interface{}{name, params},
		})
}

func ConcatKey(contract common.Address, args ...[]byte) []byte {
	temp := contract[:]
	for _, arg := range args {
		temp = append(temp, arg...)
	}
	return temp
}

func ValidateOwner(native *native.NativeService, address common.Address) error {
	if native.ContextRef.CheckWitness(address) == false {
		return errors.NewErr("validateOwner, authentication failed!")
	}
	return nil
}
