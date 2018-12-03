package main

import (
	"dan-road-vbft/account"
	"dan-road-vbft/cmd"
	cmdcom "dan-road-vbft/cmd/common"
	"dan-road-vbft/cmd/utils"
	"dan-road-vbft/common"
	"dan-road-vbft/common/config"
	"dan-road-vbft/common/log"
	"dan-road-vbft/consensus"
	"dan-road-vbft/core/genesis"
	"dan-road-vbft/core/ledger"
	"dan-road-vbft/eventbus/actor"
	"dan-road-vbft/events"
	bactor "dan-road-vbft/http/base/actor"
	hserver "dan-road-vbft/http/base/actor"
	"dan-road-vbft/p2pserver"
	netreqactor "dan-road-vbft/p2pserver/actor/req"
	p2pactor "dan-road-vbft/p2pserver/actor/server"
	"dan-road-vbft/txnpool"
	tc "dan-road-vbft/txnpool/common"
	"dan-road-vbft/txnpool/proc"
	"dan-road-vbft/validator/stateful"
	"dan-road-vbft/validator/stateless"
	"fmt"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	initApp().Run(os.Args)
}

func initApp() *cli.App {
	app := cli.NewApp()
	app.Usage = "cv-bft"
	app.Action = startApp

	app.Commands = []cli.Command{
		cmd.AccountCommand,
		cmd.AssetCommand,
		cmd.ImportCommand,
	}
	app.Flags = []cli.Flag{
		//common setting
		utils.ConfigFlag,
		utils.LogLevelFlag,
		utils.DisableEventLogFlag,
		utils.DataDirFlag,
		//account setting
		utils.WalletFileFlag,
		utils.AccountAddressFlag,
		utils.AccountPassFlag,
		//consensus setting
		utils.EnableConsensusFlag,
		utils.MaxTxInBlockFlag,
		//txpool setting
		utils.GasPriceFlag,
		utils.GasLimitFlag,
		utils.TxpoolPreExecDisableFlag,
		utils.DisableSyncVerifyTxFlag,
		utils.DisableBroadcastNetTxFlag,
		//p2p setting
		utils.ReservedPeersOnlyFlag,
		utils.ReservedPeersFileFlag,
		utils.NetworkIdFlag,
		utils.NodePortFlag,
		utils.ConsensusPortFlag,
		utils.DualPortSupportFlag,
		utils.MaxConnInBoundFlag,
		utils.MaxConnOutBoundFlag,
		utils.MaxConnInBoundForSingleIPFlag,
		//test mode setting
		utils.EnableTestModeFlag,
		utils.TestModeGenBlockTimeFlag,
		//rpc setting
		utils.RPCDisabledFlag,
		utils.RPCPortFlag,
		utils.RPCLocalEnableFlag,
		utils.RPCLocalProtFlag,
		//rest setting
		utils.RestfulEnableFlag,
		utils.RestfulPortFlag,
		//ws setting
		utils.WsEnabledFlag,
		utils.WsPortFlag,
	}

	app.Before = func(context *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return nil
	}

	return app
}

func startApp(ctx *cli.Context) {
	_, err := initConfig(ctx)
	if err != nil {
		log.Errorf("initConfig error:%s", err)
		return
	}

	acc, err := initAccount(ctx)
	if err != nil {
		log.Errorf("initWallet error:%s", err)
		return
	}
	ldg, err := initLedger(ctx)
	if err != nil {
		log.Errorf("%s", err)
		return
	}
	defer ldg.Close()
	txpool, err := initTxPool(ctx)
	if err != nil {
		log.Errorf("initTxPool error:%s", err)
		return
	}
	p2pSvr, p2pPid, err := initP2PNode(ctx, txpool)
	if err != nil {
		log.Errorf("initP2PNode error:%s", err)
		return
	}
	_, err = initConsensus(ctx, p2pPid, txpool, acc)
	if err != nil {
		log.Errorf("initConsensus error:%s", err)
		return
	}
	log.Info(p2pSvr)
	waitToExit()
}

func initConfig(ctx *cli.Context) (*config.OntologyConfig, error) {
	//init ontology config from cli
	cfg, err := cmd.SetOntologyConfig(ctx)
	if err != nil {
		return nil, err
	}
	log.Infof("Config init success")
	return cfg, nil
}

func initAccount(ctx *cli.Context) (*account.Account, error) {
	if !config.DefConfig.Consensus.EnableConsensus {
		return nil, nil
	}
	walletFile := ctx.GlobalString(utils.GetFlagName(utils.WalletFileFlag))
	if walletFile == "" {
		return nil, fmt.Errorf("Please config wallet file using --wallet flag")
	}
	if !common.FileExisted(walletFile) {
		return nil, fmt.Errorf("Cannot find wallet file:%s. Please create wallet first", walletFile)
	}
	acc, err := cmdcom.GetAccount(ctx)
	return acc, err
}

func initTxPool(ctx *cli.Context) (*proc.TXPoolServer, error) {
	disablePreExec := ctx.GlobalBool(utils.GetFlagName(utils.TxpoolPreExecDisableFlag))
	bactor.DisableSyncVerifyTx = ctx.GlobalBool(utils.GetFlagName(utils.DisableSyncVerifyTxFlag))
	disableBroadcastNetTx := ctx.GlobalBool(utils.GetFlagName(utils.DisableBroadcastNetTxFlag))
	txPoolServer, err := txnpool.StartTxnPoolServer(disablePreExec, disableBroadcastNetTx)
	if err != nil {
		return nil, fmt.Errorf("Init txpool error:%s", err)
	}
	stlValidator, _ := stateless.NewValidator("stateless_validator")
	stlValidator.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stlValidator2, _ := stateless.NewValidator("stateless_validator2")
	stlValidator2.Register(txPoolServer.GetPID(tc.VerifyRspActor))
	stfValidator, _ := stateful.NewValidator("stateful_validator")
	stfValidator.Register(txPoolServer.GetPID(tc.VerifyRspActor))

	hserver.SetTxnPoolPid(txPoolServer.GetPID(tc.TxPoolActor))
	hserver.SetTxPid(txPoolServer.GetPID(tc.TxActor))

	log.Infof("TxPool init success")
	return txPoolServer, nil
}

func initP2PNode(ctx *cli.Context, txpoolSvr *proc.TXPoolServer) (*p2pserver.P2PServer, *actor.PID, error) {

	p2p := p2pserver.NewServer()

	p2pActor := p2pactor.NewP2PActor(p2p)
	p2pPID, err := p2pActor.Start()
	if err != nil {
		return nil, nil, fmt.Errorf("p2pActor init error %s", err)
	}
	p2p.SetPID(p2pPID)
	err = p2p.Start()
	if err != nil {
		return nil, nil, fmt.Errorf("p2p service start error %s", err)
	}
	netreqactor.SetTxnPoolPid(txpoolSvr.GetPID(tc.TxActor))
	txpoolSvr.RegisterActor(tc.NetActor, p2pPID)
	hserver.SetNetServerPID(p2pPID)
	p2p.WaitForPeersStart()
	log.Infof("P2P init success")
	return p2p, p2pPID, nil
}

func initConsensus(ctx *cli.Context, p2pPid *actor.PID, txpoolSvr *proc.TXPoolServer, acc *account.Account) (consensus.ConsensusService, error) {
	pool := txpoolSvr.GetPID(tc.TxPoolActor)

	consensusService, err := consensus.NewConsensusService(acc, pool, p2pPid)

	if err != nil {
		return nil, fmt.Errorf("NewConsensusService error:%s", err)
	}
	consensusService.Start()

	netreqactor.SetConsensusPid(consensusService.GetPID())
	hserver.SetConsensusPid(consensusService.GetPID())

	log.Infof("Consensus init success")
	return consensusService, nil
}

func initLedger(ctx *cli.Context) (*ledger.Ledger, error) {
	events.Init() //Init event hub

	var err error
	dbDir := utils.GetStoreDirPath(config.DefConfig.Common.DataDir, config.DefConfig.P2PNode.NetworkName)
	//dbDir := "./chain/dan"
	ledger.DefLedger, err = ledger.NewLedger(dbDir)
	if err != nil {
		return nil, fmt.Errorf("NewLedger error:%s", err)
	}
	bookKeepers, err := config.DefConfig.GetBookkeepers()
	if err != nil {
		return nil, fmt.Errorf("GetBookkeepers error:%s", err)
	}
	genesisConfig := config.DefConfig.Genesis
	genesisBlock, err := genesis.BuildGenesisBlock(bookKeepers, genesisConfig)
	if err != nil {
		return nil, fmt.Errorf("genesisBlock error %s", err)
	}
	err = ledger.DefLedger.Init(bookKeepers, genesisBlock)
	if err != nil {
		return nil, fmt.Errorf("Init ledger error:%s", err)
	}

	log.Infof("Ledger init success")
	return ledger.DefLedger, nil
}

func waitToExit() {
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			log.Infof("Ontology received exit signal:%v.", sig.String())
			close(exit)
			break
		}
	}()
	<-exit
}
