package generate

import log "github.com/Futuremine-chain/futuremine/log/log15"

const module = "generate"

type Generate struct {
}

func NewGenerate() *Generate {
	return &Generate{}
}

func (g *Generate) Name() string {
	return module
}

func (g *Generate) Start() error {
	log.Info("Generate started successfully", "module", module)
	return nil
}

func (g *Generate) Stop() error {
	return nil
}
