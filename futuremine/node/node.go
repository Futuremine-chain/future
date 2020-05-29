package node

import (
	"github.com/Futuremine-chain/futuremine/server"
	log "github.com/Futuremine-chain/futuremine/tools/log/log15"
)

const module = "fmc_node"

type FMCNode struct {
	services []server.IService
}

func NewFMCNode() *FMCNode {
	return &FMCNode{
		services: make([]server.IService, 0),
	}
}

func (fmc *FMCNode) Start() error {

	if err := fmc.startServices(); err != nil {
		return err
	}
	return nil
}

func (fmc *FMCNode) Stop() error {
	for _, s := range fmc.services {
		if err := s.Stop(); err != nil {
			log.Error("Service failed to stop", "module", module, "service", s.Name(), "error", err.Error())
		}
	}
	return nil
}

func (fmc *FMCNode) Register(s server.IService) {
	fmc.services = append(fmc.services, s)
}

func (fmc *FMCNode) startServices() error {
	for _, s := range fmc.services {
		if err := s.Start(); err != nil {
			log.Error("Service failed to start", "module", module, "service", s.Name(), "error", err.Error())
		}
	}
	return nil
}
