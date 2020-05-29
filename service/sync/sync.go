package sync

import log "github.com/Futuremine-chain/futuremine/tools/log/log15"

const module = "sync"

type Sync struct {
}

func NewSync() *Sync {
	return &Sync{}
}

func (s *Sync) Name() string {
	return module
}

func (s *Sync) Start() error {
	log.Info("Sync started successfully", "module", module)
	return nil
}

func (s *Sync) Stop() error {
	return nil
}
