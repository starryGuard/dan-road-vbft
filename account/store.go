package account

import (
	"dan-road-vbft/common"
)

type ClientStore interface {
	BuildDatabase(path string)

	SaveStoredData(name string, value []byte)

	LoadStoredData(name string) []byte

	LoadAccount() map[common.Address]*Account
}
