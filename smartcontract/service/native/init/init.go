package init

import (
	"bytes"
	"math/big"

	"dan-road-vbft/common"
	"dan-road-vbft/smartcontract/service/native/auth"
	params "dan-road-vbft/smartcontract/service/native/global_params"
	"dan-road-vbft/smartcontract/service/native/governance"
	"dan-road-vbft/smartcontract/service/native/ong"
	"dan-road-vbft/smartcontract/service/native/ont"
	"dan-road-vbft/smartcontract/service/native/ontid"
	"dan-road-vbft/smartcontract/service/native/utils"
	"dan-road-vbft/smartcontract/service/neovm"
	vm "dan-road-vbft/vm/neovm"
)

var (
	COMMIT_DPOS_BYTES = InitBytes(utils.GovernanceContractAddress, governance.COMMIT_DPOS)
)

func init() {
	ong.InitOng()
	ont.InitOnt()
	params.InitGlobalParams()
	ontid.Init()
	auth.Init()
	governance.InitGovernance()
}

func InitBytes(addr common.Address, method string) []byte {
	bf := new(bytes.Buffer)
	builder := vm.NewParamsBuilder(bf)
	builder.EmitPushByteArray([]byte{})
	builder.EmitPushByteArray([]byte(method))
	builder.EmitPushByteArray(addr[:])
	builder.EmitPushInteger(big.NewInt(0))
	builder.Emit(vm.SYSCALL)
	builder.EmitPushByteArray([]byte(neovm.NATIVE_INVOKE_NAME))

	return builder.ToArray()
}
