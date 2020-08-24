package main

import (
	"fmt"
	"github.com/Futuremine-chain/future/common/config"
	"github.com/Futuremine-chain/future/common/horn"
	"github.com/Futuremine-chain/future/future/common/blockchain"
	fmcdpos "github.com/Futuremine-chain/future/future/common/dpos"
	"github.com/Futuremine-chain/future/future/common/msglist"
	"github.com/Futuremine-chain/future/future/common/private"
	fmcstatus "github.com/Futuremine-chain/future/future/common/status"
	"github.com/Futuremine-chain/future/future/common/status/act_status"
	"github.com/Futuremine-chain/future/future/common/status/dpos_status"
	"github.com/Futuremine-chain/future/future/common/status/token_status"
	"github.com/Futuremine-chain/future/future/node"
	"github.com/Futuremine-chain/future/future/request"
	"github.com/Futuremine-chain/future/future/rpc"
	"github.com/Futuremine-chain/future/service/generate"
	"github.com/Futuremine-chain/future/service/gorutinue"
	"github.com/Futuremine-chain/future/service/p2p"
	"github.com/Futuremine-chain/future/service/peers"
	"github.com/Futuremine-chain/future/service/pool"
	sync_service "github.com/Futuremine-chain/future/service/sync"
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
	if err := FCMain(); err != nil {
		fmt.Println("Failed to start, ", err)
		os.Exit(1)
	}
}

// main start the FMC node function
func FCMain() error {
	var node *node.FCNode
	var err error
	wg := sync.WaitGroup{}
	wg.Add(1)

	if node, err = createFCNode(); err != nil {
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

func createFCNode() (*node.FCNode, error) {
	err := config.LoadParam(private.NewPrivate(nil))
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

	dPos := fmcdpos.NewDPos(dPosStatus)
	status := fmcstatus.NewFMCStatus(actStatus, dPosStatus, tokenStatus)
	gPool := gorutinue.NewPool()
	chain, err := blockchain.NewFMCChain(status, dPos)
	if err != nil {
		return nil, err
	}
	reqHandler := request.NewRequestHandler(chain)
	peersSv := peers.NewPeers(reqHandler)

	p2pSv, err := p2p.NewP2p(peersSv, reqHandler)
	if err != nil {
		return nil, err
	}

	horn := horn.NewHorn(peersSv, gPool, reqHandler)
	msgManage, err := msglist.NewMsgManagement(status, actStatus)
	if err != nil {
		return nil, err
	}
	poolSv := pool.NewPool(horn, msgManage)

	rpcSv := rpc.NewRpc(status, poolSv, chain, peersSv)
	syncSv := sync_service.NewSync(peersSv, dPosStatus, reqHandler, chain)
	generateSv := generate.NewGenerate(chain, dPos, poolSv, horn)
	node := node.NewFCNode()

	rpcSv.RegisterLocalInfo(node.LocalInfo)
	reqHandler.RegisterLocalInfo(node.LocalInfo)

	chain.RegisterMsgPoolDeleteFunc(poolSv.Delete)

	// Register peer nodes to send blocks and message processing
	reqHandler.RegisterReceiveMessage(poolSv.ReceiveMsgFromPeer)
	reqHandler.RegisterReceiveBlock(syncSv.ReceivedBlockFromPeer)

	node.Register(syncSv)
	node.Register(peersSv)
	node.Register(p2pSv)
	node.Register(rpcSv)
	node.Register(reqHandler)
	node.Register(gPool)
	node.Register(poolSv)
	node.Register(generateSv)
	return node, nil
}
