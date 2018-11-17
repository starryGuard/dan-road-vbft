package event

import (
	"dan-road-vbft/common"
)

// LogEventArgs describe smart contract event log struct
type LogEventArgs struct {
	TxHash          common.Uint256
	ContractAddress common.Address
	Message         string
}
