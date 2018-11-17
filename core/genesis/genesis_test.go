package genesis

import (
	"dan-road-vbft/common"
	"dan-road-vbft/common/config"
	"dan-road-vbft/common/log"
	"dan-road-vbft/crypto/keypair"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.InitLog(0, log.Stdout)
	m.Run()
	os.RemoveAll("./ActorLog")
}

func TestGenesisBlockInit(t *testing.T) {
	_, pub, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	conf := &config.GenesisConfig{}
	block, err := BuildGenesisBlock([]keypair.PublicKey{pub}, conf)
	assert.Nil(t, err)
	assert.NotNil(t, block)
	assert.NotEqual(t, block.Header.TransactionsRoot, common.UINT256_EMPTY)
}

func TestNewParamDeployAndInit(t *testing.T) {
	deployTx := newParamContract()
	initTx := newParamInit()
	assert.NotNil(t, deployTx)
	assert.NotNil(t, initTx)
}
