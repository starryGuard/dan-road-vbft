package ong

import (
	"math/big"

	"dan-road-vbft/common"
	"dan-road-vbft/common/constants"
	"dan-road-vbft/errors"
	"dan-road-vbft/smartcontract/service/native"
	"dan-road-vbft/smartcontract/service/native/ont"
	"dan-road-vbft/smartcontract/service/native/utils"
	"dan-road-vbft/vm/neovm/types"
	"fmt"
)

func InitOng() {
	native.Contracts[utils.OngContractAddress] = RegisterOngContract
}

func RegisterOngContract(native *native.NativeService) {
	native.Register(ont.INIT_NAME, OngInit)
	native.Register(ont.TRANSFER_NAME, OngTransfer)
	native.Register(ont.APPROVE_NAME, OngApprove)
	native.Register(ont.TRANSFERFROM_NAME, OngTransferFrom)
	native.Register(ont.NAME_NAME, OngName)
	native.Register(ont.SYMBOL_NAME, OngSymbol)
	native.Register(ont.DECIMALS_NAME, OngDecimals)
	native.Register(ont.TOTALSUPPLY_NAME, OngTotalSupply)
	native.Register(ont.BALANCEOF_NAME, OngBalanceOf)
	native.Register(ont.ALLOWANCE_NAME, OngAllowance)
}

func OngInit(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	amount, err := utils.GetStorageUInt64(native, ont.GenTotalSupplyKey(contract))
	if err != nil {
		return utils.BYTE_FALSE, err
	}

	if amount > 0 {
		return utils.BYTE_FALSE, errors.NewErr("Init ong has been completed!")
	}

	item := utils.GenUInt64StorageItem(constants.ONG_TOTAL_SUPPLY)
	native.CacheDB.Put(ont.GenTotalSupplyKey(contract), item.ToArray())
	native.CacheDB.Put(append(contract[:], utils.OntContractAddress[:]...), item.ToArray())
	ont.AddNotifications(native, contract, &ont.State{To: utils.OntContractAddress, Value: constants.ONG_TOTAL_SUPPLY})
	return utils.BYTE_TRUE, nil
}

func OngTransfer(native *native.NativeService) ([]byte, error) {
	var transfers ont.Transfers
	source := common.NewZeroCopySource(native.Input)
	if err := transfers.Deserialization(source); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[OngTransfer] Transfers deserialize error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	for _, v := range transfers.States {
		if v.Value == 0 {
			continue
		}
		if v.Value > constants.ONG_TOTAL_SUPPLY {
			return utils.BYTE_FALSE, fmt.Errorf("transfer ong amount:%d over totalSupply:%d", v.Value, constants.ONG_TOTAL_SUPPLY)
		}
		if _, _, err := ont.Transfer(native, contract, &v); err != nil {
			return utils.BYTE_FALSE, err
		}
		ont.AddNotifications(native, contract, &v)
	}
	return utils.BYTE_TRUE, nil
}

func OngApprove(native *native.NativeService) ([]byte, error) {
	var state ont.State
	source := common.NewZeroCopySource(native.Input)
	if err := state.Deserialization(source); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[OngApprove] state deserialize error!")
	}
	if state.Value == 0 {
		return utils.BYTE_FALSE, nil
	}
	if state.Value > constants.ONG_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("approve ong amount:%d over totalSupply:%d", state.Value, constants.ONG_TOTAL_SUPPLY)
	}
	if native.ContextRef.CheckWitness(state.From) == false {
		return utils.BYTE_FALSE, errors.NewErr("authentication failed!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	native.CacheDB.Put(ont.GenApproveKey(contract, state.From, state.To), utils.GenUInt64StorageItem(state.Value).ToArray())
	return utils.BYTE_TRUE, nil
}

func OngTransferFrom(native *native.NativeService) ([]byte, error) {
	var state ont.TransferFrom
	source := common.NewZeroCopySource(native.Input)
	if err := state.Deserialization(source); err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[OntTransferFrom] State deserialize error!")
	}
	if state.Value == 0 {
		return utils.BYTE_FALSE, nil
	}
	if state.Value > constants.ONG_TOTAL_SUPPLY {
		return utils.BYTE_FALSE, fmt.Errorf("approve ong amount:%d over totalSupply:%d", state.Value, constants.ONG_TOTAL_SUPPLY)
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	if _, _, err := ont.TransferedFrom(native, contract, &state); err != nil {
		return utils.BYTE_FALSE, err
	}
	ont.AddNotifications(native, contract, &ont.State{From: state.From, To: state.To, Value: state.Value})
	return utils.BYTE_TRUE, nil
}

func OngName(native *native.NativeService) ([]byte, error) {
	return []byte(constants.ONG_NAME), nil
}

func OngDecimals(native *native.NativeService) ([]byte, error) {
	return big.NewInt(int64(constants.ONG_DECIMALS)).Bytes(), nil
}

func OngSymbol(native *native.NativeService) ([]byte, error) {
	return []byte(constants.ONG_SYMBOL), nil
}

func OngTotalSupply(native *native.NativeService) ([]byte, error) {
	contract := native.ContextRef.CurrentContext().ContractAddress
	amount, err := utils.GetStorageUInt64(native, ont.GenTotalSupplyKey(contract))
	if err != nil {
		return utils.BYTE_FALSE, errors.NewDetailErr(err, errors.ErrNoCode, "[OntTotalSupply] get totalSupply error!")
	}
	return types.BigIntToBytes(big.NewInt(int64(amount))), nil
}

func OngBalanceOf(native *native.NativeService) ([]byte, error) {
	return ont.GetBalanceValue(native, ont.TRANSFER_FLAG)
}

func OngAllowance(native *native.NativeService) ([]byte, error) {
	return ont.GetBalanceValue(native, ont.APPROVE_FLAG)
}
