package vote

import (
	"dan-road-vbft/core/genesis"
	"dan-road-vbft/core/types"
	"dan-road-vbft/crypto/keypair"
)

func GetValidators(txs []*types.Transaction) ([]keypair.PublicKey, error) {
	// TODO implement vote
	return genesis.GenesisBookkeepers, nil
}
