package main

import (
	"fmt"
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/common/horn"
	"github.com/Futuremine-chain/futuremine/futuremine/common/blockchain"
	dpos2 "github.com/Futuremine-chain/futuremine/futuremine/common/dpos"
	"github.com/Futuremine-chain/futuremine/futuremine/common/private"
	status2 "github.com/Futuremine-chain/futuremine/futuremine/common/status"
	"github.com/Futuremine-chain/futuremine/futuremine/common/status/act_status"
	"github.com/Futuremine-chain/futuremine/futuremine/common/status/dpos_status"
	"github.com/Futuremine-chain/futuremine/futuremine/common/status/token_status"
	"github.com/Futuremine-chain/futuremine/futuremine/common/txlist"
	"github.com/Futuremine-chain/futuremine/futuremine/node"
	"github.com/Futuremine-chain/futuremine/futuremine/request"
	"github.com/Futuremine-chain/futuremine/service/generate"
	"github.com/Futuremine-chain/futuremine/service/gorutinue"
	"github.com/Futuremine-chain/futuremine/service/p2p"
	"github.com/Futuremine-chain/futuremine/service/peers"
	"github.com/Futuremine-chain/futuremine/service/pool"
	"github.com/Futuremine-chain/futuremine/service/rpc"
	sync_service "github.com/Futuremine-chain/futuremine/service/sync"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

// interruptSignals defines the default signals to catch in order to do a proper
// shutdown.  This may be modified during init depending on the platform.
var interruptSignals = []os.Signal{
	os.Interrupt,
	os.Kill,
	syscall.SIGINT,
	syscall.SIGTERM,
}

func main() {
	// Initialize the goroutine count,  Use all processor cores.
	runtime.GOMAXPROCS(runtime.NumCPU())
	// Work around defer not working after os.Exit()
	if err := FMCMain(); err != nil {
		fmt.Println("faild to start, ", err)
		os.Exit(1)
	}
}

// main start the FMC node function
func FMCMain() error {
	var node *node.FMCNode
	var err error
	wg := sync.WaitGroup{}
	wg.Add(1)

	if node, err = createFMCNode(); err != nil {
		return err
	}
	if err := node.Start(); err != nil {
		return err
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, interruptSignals...)

	// Listen for initial shutdown signal and close the returned
	// channel to notify the caller.
	go func() {
		<-c
		node.Stop()
		close(c)
		wg.Done()
	}()
	wg.Wait()
	return nil
}

func createFMCNode() (*node.FMCNode, error) {
	cfg, err := config.LoadConfig(NewFmcApp(), private.NewPrivate())
	if err != nil {
		return nil, err
	}
	actStatus, err := act_status.NewActStatus()
	if err != nil {
		return nil, err
	}
	dPosStatus, err := dpos_status.NewDPosStatus()
	if err != nil {
		return nil, err
	}
	tokenStatus, err := token_status.NewTokenStatus()
	if err != nil {
		return nil, err
	}

	dpos := dpos2.NewDPos()
	status := status2.NewFMCStatus(actStatus, dPosStatus, tokenStatus)
	chain, err := blockchain.NewFMCChain(status, dpos)
	if err != nil {
		return nil, err
	}
	peersSv := peers.NewPeers()
	reqHandler := request.NewRequestHandler(chain)
	p2pSv, err := p2p.NewP2p(cfg, peersSv, reqHandler)
	if err != nil {
		return nil, err
	}
	txmanage, err := txlist.NewTxManagement(nil, nil)
	if err != nil {
		return nil, err
	}
	rpcSv := rpc.NewRpc()
	gPool := gorutinue.NewPool()
	horn := horn.NewHorn(peersSv, gPool)
	poolSv := pool.NewPool(horn, txmanage)
	syncSv := sync_service.NewSync(peersSv, dpos, reqHandler, chain)
	generateSv := generate.NewGenerate(chain)
	node := node.NewFMCNode()

	node.Register(peersSv)
	node.Register(p2pSv)
	node.Register(rpcSv)
	node.Register(reqHandler)
	node.Register(gPool)
	node.Register(poolSv)
	node.Register(syncSv)
	node.Register(generateSv)
	return node, nil
}
