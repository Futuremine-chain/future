package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Futuremine-chain/futuremine/command/wallet/command"
	"github.com/Futuremine-chain/futuremine/command/wallet/config"
	"github.com/Futuremine-chain/futuremine/futuremine/common/param"
	"github.com/Futuremine-chain/futuremine/tools/utils"
	"github.com/spf13/cobra"
	"os"
)

var preConfig *config.Config
var (
	defaultFormat      = true
	defaultTestNet     = false
	defaultKeyStoreDir = "keystore"
	defaultRpcCer      = "server.pem"
	defaultRpcIp       = "127.0.0.1"
)

func init() {
	cobra.EnableCommandSorting = true
	bindFlags()
}

func main() {
	command.RootCmd.PersistentPreRunE = LoadConfig
	if err := command.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func bindFlags() {
	preConfig = &config.Config{}
	gFlags := command.RootCmd.PersistentFlags()

	gFlags.StringVarP(&preConfig.ConfigFile, "config", "c", "wallet.toml", "Wallet profile")
}

// LoadConfig config file and flags
func LoadConfig(cmd *cobra.Command, args []string) (err error) {
	fileCfg := &config.Config{}
	_, err = toml.DecodeFile(preConfig.ConfigFile, fileCfg)
	if err != nil {
		if !cmd.Flag("config").Changed {
			if fExit := utils.Exist(preConfig.ConfigFile); !fExit {
				return fmt.Errorf("config file is not exist")
			}
			_, err = toml.DecodeFile(cmd.Flag("config").Value.String(), fileCfg)
			if err != nil {
				return fmt.Errorf("config file %s is not exist", cmd.Flag("config").Value.String())
			}
		}
	}
	if fileCfg.KeystoreDir == "" {
		fileCfg.KeystoreDir = defaultKeyStoreDir
	}

	if fileCfg.RpcPort == "" {
		fileCfg.RpcPort = "19161"
	}

	if fileCfg.RpcIp == "" {
		fileCfg.RpcIp = defaultRpcIp
	}

	command.Cfg = fileCfg
	if command.Cfg.TestNet {
		command.Net = param.TestNet
	}
	return nil
}
