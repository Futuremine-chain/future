package peers

import log "github.com/Futuremine-chain/futuremine/log/log15"

const module = "peers"

type Peers struct {
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
