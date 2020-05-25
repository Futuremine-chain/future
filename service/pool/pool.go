package pool

import log "github.com/Futuremine-chain/futuremine/log/log15"

const module = "pool"

type Pool struct {
}

func NewPool() *Pool {
	return &Pool{}
}

func (p *Pool) Name() string {
	return "pool"
}

func (p *Pool) Start() error {
	log.Info("Pool started successfully", "module", module)
	return nil
}

func (p *Pool) Stop() error {
	return nil
}
