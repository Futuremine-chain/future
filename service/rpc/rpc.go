package rpc

import log "github.com/Futuremine-chain/futuremine/tools/log/log15"

const module = "rpc"

type Rpc struct {
}

func NewRpc() *Rpc {
	return &Rpc{}
}

func (r *Rpc) Name() string {
	return module
}

func (r *Rpc) Start() error {
	log.Info("Rpc started successfully", "module", module)
	return nil
}

func (r *Rpc) Stop() error {
	return nil
}
