package generate

import (
	"github.com/Futuremine-chain/futuremine/common/blockchain"
	"github.com/Futuremine-chain/futuremine/common/dpos"
	"github.com/Futuremine-chain/futuremine/common/horn"
	"github.com/Futuremine-chain/futuremine/service/pool"
	log "github.com/Futuremine-chain/futuremine/tools/log/log15"
	"time"
)

const (
	module           = "generate"
	maxPackedTxCount = 999
)

type Generate struct {
	horn        *horn.Horn
	dPos        dpos.IDPos
	chain       blockchain.IBlockChain
	pool        pool.Pool
	minerWorkCh chan bool
	stop        chan bool
	stopped     chan bool
}

func NewGenerate(chain blockchain.IBlockChain) *Generate {
	return &Generate{
		chain:   chain,
		stop:    make(chan bool),
		stopped: make(chan bool),
	}
}

func (g *Generate) Name() string {
	return module
}

func (g *Generate) Start() error {
	log.Info("Generate started successfully", "module", module)
	return nil
}

func (g *Generate) Stop() error {
	return nil
}

func (g *Generate) Generate() {
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case _, _ = <-g.stop:
			g.stopped <- true
			log.Info("Stop generate block")
			return
		case _ = <-ticker:
			g.generateBlock()
		}
	}
}

func (g *Generate) generateBlock() {
	txs := g.pool.NeedPackaged(maxPackedTxCount)
	nextBlock := g.chain.NextBlock(txs)
	// Check if it is your turn to make blocks
	err := g.dPos.CheckSigner(g.chain, nextBlock)
	if err != nil {
		//.Warn("check winner failed!", "height", header.Height, "error", err)
		return
	}
	log.Info("Block generation successful", "module", module,
		"height", nextBlock.Hash().String(),
		"hash", nextBlock.Hash().String(),
		"signer", nextBlock.Signer().String(),
	)
	g.horn.BroadcastBlock(nextBlock)
}
