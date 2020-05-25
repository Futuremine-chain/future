package p2p

import log "github.com/Futuremine-chain/futuremine/log/log15"

const module = "p2p"

type P2p struct {
}

func NewP2p() *P2p {
	return &P2p{}
}

func (p *P2p) Name() string {
	return "p2p"
}

func (p *P2p) Start() error {
	log.Info("P2P started successfully", "module", module)
	return nil
}

func (p *P2p) Stop() error {
	return nil
}
