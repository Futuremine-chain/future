package node

import (
	"encoding/json"
	"github.com/Futuremine-chain/future/common/config"
	"github.com/Futuremine-chain/future/common/param"
	"github.com/Futuremine-chain/future/server"
	log "github.com/Futuremine-chain/future/tools/log/log15"
	"github.com/Futuremine-chain/future/types"
)

const module = "fc_node"

type FCNode struct {
	services []server.IService
}

func NewFCNode() *FCNode {
	return &FCNode{
		services: make([]server.IService, 0),
	}
}

func (fmc *FCNode) Start() error {
	if err := fmc.startServices(); err != nil {
		return err
	}
	return nil
}

func (fmc *FCNode) Stop() error {
	for _, s := range fmc.services {
		if err := s.Stop(); err != nil {
			log.Error("Service failed to stop", "module", module, "service", s.Name(), "error", err.Error())
		}
	}

	return nil
}

func (fmc *FCNode) Register(s server.IService) {
	fmc.services = append(fmc.services, s)
}

func (fmc *FCNode) LocalInfo() *types.Local {
	all := make(map[string]interface{})
	for _, s := range fmc.services {
		infoMap := s.Info()
		for name, value := range infoMap {
			all[name] = value
		}
	}
	all["version"] = param.Version
	all["network"] = config.Param.Name
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

func (fmc *FCNode) startServices() error {
	for _, s := range fmc.services {
		if err := s.Start(); err != nil {
			log.Error("Service failed to start", "module", module, "service", s.Name(), "error", err.Error())
		}
	}
	return nil
}
