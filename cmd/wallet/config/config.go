package config

import (
	"github.com/Futuremine-chain/future/common/config"
)

type Config struct {
	ConfigFile  string
	Format      bool
	TestNet     bool
	KeystoreDir string
	config.RpcConfig
}
