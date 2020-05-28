package peers

import (
	log "github.com/Futuremine-chain/futuremine/log/log15"
	"github.com/libp2p/go-libp2p-core/peer"
)

const module = "peers"

type Peers struct {
	PeerExist func(info *peer.AddrInfo) bool
	Add       func(*Peer)
	Remove    func(string)
}

func NewPeers() *Peers {
	return &Peers{}
}

func (p *Peers) Name() string {
	return "peers"
}

func (p *Peers) Start() error {
	log.Info("Peers started successfully", "module", module)
	return nil
}

func (p *Peers) Stop() error {
	return nil
}
