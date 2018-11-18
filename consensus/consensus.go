package consensus

import (
	"dan-road-vbft/account"
	"dan-road-vbft/consensus/vbft"
	"dan-road-vbft/eventbus/actor"
)

type ConsensusService interface {
	Start() error
	Halt() error
	GetPID() *actor.PID
}

func NewConsensusService(account *account.Account, txpool *actor.PID, p2p *actor.PID) (ConsensusService, error) {
	var consensus ConsensusService
	var err error
	consensus, err = vbft.NewVbftServer(account, txpool, p2p)
	return consensus, err
}
