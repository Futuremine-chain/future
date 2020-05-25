package main

import (
	"github.com/Futuremine-chain/futuremine/fmc/node"
	"github.com/Futuremine-chain/futuremine/service/connect"
	"github.com/Futuremine-chain/futuremine/service/generate"
	"github.com/Futuremine-chain/futuremine/service/p2p"
	"github.com/Futuremine-chain/futuremine/service/peers"
	"github.com/Futuremine-chain/futuremine/service/pool"
	"github.com/Futuremine-chain/futuremine/service/rpc"
	"github.com/Futuremine-chain/futuremine/service/sync"
	"os"
	"os/signal"
	"runtime"
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
	peersSv := peers.NewPeers()
	p2pSv := p2p.NewP2p()
	rpcSv := rpc.NewRpc()
	connectSv := connect.NewConnect()
	poolSv := pool.NewPool()
	syncSv := sync.NewSync()
	generateSv := generate.NewGenerate()
	node := node.NewFMCNode()
	node.Register(peersSv)
	node.Register(p2pSv)
	node.Register(rpcSv)
	node.Register(connectSv)
	node.Register(poolSv)
	node.Register(syncSv)
	node.Register(generateSv)
	return node, nil
}