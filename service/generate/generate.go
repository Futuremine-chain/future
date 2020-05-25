package generate

import log "github.com/Futuremine-chain/futuremine/log/log15"

const module = "pool"

type Generate struct {
}

func NewGenerate() *Generate {
	return &Generate{}
}

func (g *Generate) Name() string {
	return "generate"
}

func (g *Generate) Start() error {
	log.Info("Generate started successfully", "module", module)
	return nil
}

func (g *Generate) Stop() error {
	return nil
}
