package operators

import (
	"dan-road-vbft/vm/wasmvm/wasm"
)

var (
	I32Const = newOp(0x41, "i32.const", nil, wasm.ValueTypeI32)
	I64Const = newOp(0x42, "i64.const", nil, wasm.ValueTypeI64)
	F32Const = newOp(0x43, "f32.const", nil, wasm.ValueTypeF32)
	F64Const = newOp(0x44, "f64.const", nil, wasm.ValueTypeF64)
)
