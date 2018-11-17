package req

import (
	"dan-road-vbft/eventbus/actor"
)

var ConsensusPid *actor.PID

func SetConsensusPid(conPid *actor.PID) {
	ConsensusPid = conPid
}
