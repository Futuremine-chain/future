package request

import (
	log "github.com/Futuremine-chain/futuremine/log/log15"
	core "github.com/libp2p/go-libp2p-core"
)

const module = "request"

type RequestHandler struct {
	HandleRequest func(stream core.Stream)
}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{}
}

func (c *RequestHandler) Name() string {
	return "RequestHandler"
}

func (c *RequestHandler) Start() error {
	log.Info("Request handler started successfully", "module", module)
	return nil
}

func (c *RequestHandler) Stop() error {
	return nil
}
