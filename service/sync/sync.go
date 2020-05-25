package sync

import log "github.com/Futuremine-chain/futuremine/log/log15"

const module = "sync"

type Sync struct {
}

func NewSync() *Sync {
	return &Sync{}
}

func (s *Sync) Name() string {
	return "sync"
}

func (s *Sync) Start() error {
	log.Info("Sync started successfully", "module", module)
	return nil
}

func (s *Sync) Stop() error {
	return nil
}
