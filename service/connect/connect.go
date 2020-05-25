package connect

import log "github.com/Futuremine-chain/futuremine/log/log15"

const module = "connect"

type Connect struct {
}

func NewConnect() *Connect {
	return &Connect{}
}

func (c *Connect) Name() string {
	return "connect"
}

func (c *Connect) Start() error {
	log.Info("Connect started successfully", "module", module)
	return nil
}

func (c *Connect) Stop() error {
	return nil
}
