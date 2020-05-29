package pool

import (
	"fmt"
	"github.com/Futuremine-chain/futuremine/common/horn"
	"github.com/Futuremine-chain/futuremine/common/txlist"
	log "github.com/Futuremine-chain/futuremine/log/log15"
	"github.com/Futuremine-chain/futuremine/tools/utils"
	"github.com/Futuremine-chain/futuremine/types"
	"time"
)

const module = "pool"

// Clear the expired transaction interval
const (
	txExpiredTime     = 60 * 60 * 3
	monitorTxInterval = 2
	maxPoolTx         = 5000
)

type Pool struct {
	txMgt       txlist.ITxList
	horn        *horn.Horn
	broadcastCh chan types.ITransaction
}

func NewPool(horn *horn.Horn, txMgt txlist.ITxList) *Pool {
	return &Pool{
		horn:        horn,
		broadcastCh: make(chan types.ITransaction, 100),
	}
}

func (p *Pool) Name() string {
	return module
}

func (p *Pool) Start() error {
	if err := p.txMgt.Read(); err != nil {
		log.Error("The transaction pool failed to read the transaction", "module", module, "error", err)
		return err
	}
	go p.monitorExpired()
	go p.startChan()
	log.Info("Pool started successfully", "module", module)
	return nil
}

func (p *Pool) Stop() error {
	if err := p.txMgt.Close(); err != nil {
		return err
	}
	log.Info("Transaction pool stopped", "module", module)
	return nil
}

// Verify adding transactions to the transaction pool
func (p *Pool) Put(tx types.ITransaction, isPeer bool) error {
	if err := p.txMgt.Put(tx); err != nil {
		return utils.Error(fmt.Sprintf("add transaction failed, %s", err.Error()), module)
	}
	log.Info("Received the transaction", "module", module, "hash", tx.Hash().String())
	if !isPeer {
		p.broadcastCh <- tx
	}
	return nil
}

func (p *Pool) startChan() {
	for {
		select {
		case tx := <-p.broadcastCh:
			p.horn.BroadcastTx(tx)
		}
	}
}

func (p *Pool) monitorExpired() {
	t := time.NewTicker(time.Second * monitorTxInterval)
	defer t.Stop()

	for range t.C {
		p.removeExpired()
	}
}

func (p *Pool) removeExpired() {
	threshold := utils.NowUnix() - txExpiredTime
	p.txMgt.DeleteExpired(threshold)
}

func (p *Pool) DeleteAndUpdate(txs types.ITransactions) {
	p.txMgt.DeleteAndUpdate(txs)
}
