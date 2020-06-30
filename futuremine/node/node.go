package node

import (
	"encoding/json"
	"github.com/Futuremine-chain/futuremine/server"
	log "github.com/Futuremine-chain/futuremine/tools/log/log15"
	"github.com/Futuremine-chain/futuremine/types"
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

func (fmc *FMCNode) LocalInfo() *types.Local {
	all := make(map[string]interface{})
	for _, s := range fmc.services {
		infoMap := s.Info()
		for name, value := range infoMap {
			all[name] = value
		}
	}
	bytes, err := json.Marshal(all)
	if err != nil {
		return &types.Local{}
	}
	var rs *types.Local
	err = json.Unmarshal(bytes, &rs)
	if err != nil {
		return &types.Local{}
	}
	return rs
}

func (fmc *FMCNode) startServices() error {
	for _, s := range fmc.services {
		if err := s.Start(); err != nil {
			log.Error("Service failed to start", "module", module, "service", s.Name(), "error", err.Error())
		}
	}
	return nil
}
