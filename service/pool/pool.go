package pool

import (
	"fmt"
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/common/horn"
	"github.com/Futuremine-chain/futuremine/common/msglist"
	log "github.com/Futuremine-chain/futuremine/tools/log/log15"
	"github.com/Futuremine-chain/futuremine/tools/utils"
	"github.com/Futuremine-chain/futuremine/types"
	"time"
)

const module = "pool"

type Pool struct {
	msgMgt      msglist.IMsgList
	horn        *horn.Horn
	broadcastCh chan types.IMessage
	deleteMsg   chan types.IMessage
	close       chan bool
}

func NewPool(horn *horn.Horn, msgMgt msglist.IMsgList) *Pool {
	pool := &Pool{
		msgMgt:      msgMgt,
		horn:        horn,
		broadcastCh: make(chan types.IMessage, 100),
		deleteMsg:   make(chan types.IMessage, 10000),
		close:       make(chan bool),
	}
	return pool
}

func (p *Pool) Name() string {
	return module
}

func (p *Pool) Start() error {
	if err := p.msgMgt.Read(); err != nil {
		log.Error("The message pool failed to read the message", "module", module, "error", err)
		return err
	}
	go p.monitorExpired()
	go p.startChan()
	log.Info("Pool started successfully", "module", module)
	return nil
}

func (p *Pool) Stop() error {
	if err := p.msgMgt.Close(); err != nil {
		p.close <- true
		return err
	}
	p.close <- true
	log.Info("Message pool was stopped", "module", module)
	return nil
}

// Verify adding messages to the message pool
func (p *Pool) Put(msg types.IMessage, isPeer bool) error {
	if err := p.msgMgt.Put(msg); err != nil {
		return utils.Error(fmt.Sprintf("add message failed, %s", err.Error()), module)
	}
	log.Info("Received the message", "module", module, "hash", msg.Hash().String())
	if !isPeer {
		p.broadcastCh <- msg
	}
	return nil
}

func (p *Pool) NeedPackaged(count int) []types.IMessage {
	msgs := p.msgMgt.NeedPackaged(count)
	return msgs
}

func (p *Pool) startChan() {
	for {
		select {
		case _ = <-p.close:
			return
		case msg := <-p.broadcastCh:
			p.horn.BroadcastMsg(msg)
		case msg := <-p.deleteMsg:
			p.msgMgt.Delete(msg)
		}
	}
}

func (p *Pool) ReceiveMsgFromPeer(msg types.IMessage) error {
	return p.Put(msg, true)
}

func (p *Pool) monitorExpired() {
	t := time.NewTicker(time.Second * config.Param.MonitorMsgInterval)
	defer t.Stop()

	for range t.C {
		p.removeExpired()
	}
}

func (p *Pool) removeExpired() {
	threshold := utils.NowUnix() - config.Param.MsgExpiredTime
	p.msgMgt.DeleteExpired(threshold)
}

func (p *Pool) Delete(msg types.IMessage) {
	p.deleteMsg <- msg
}

// Get all transactions in the trading pool
func (p *Pool) All() ([]types.IMessage, []types.IMessage) {
	prepareTxs, futureTxs := p.msgMgt.GetAll()
	return prepareTxs, futureTxs
}
