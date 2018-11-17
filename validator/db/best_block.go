package db

import (
	"dan-road-vbft/common"
	"dan-road-vbft/core/types"
)

type BestBlock struct {
	Height uint32
	Hash   common.Uint256
}

type BestStateProvider interface {
	GetBestBlock() (BestBlock, error)
	GetBestHeader() (*types.Header, error)
}
