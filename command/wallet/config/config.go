package config

import (
	"github.com/Futuremine-chain/futuremine/common/config"
)

type Config struct {
	ConfigFile  string
	Format      bool
	TestNet     bool
	KeystoreDir string
	config.RpcConfig
}
