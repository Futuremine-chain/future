package config

import (
	"github.com/BurntSushi/toml"
	log2 "github.com/Futuremine-chain/futuremine/log"
	log "github.com/Futuremine-chain/futuremine/log/log15"
	"github.com/Futuremine-chain/futuremine/tools/utils"
	"github.com/jessevdk/go-flags"
	"os"
	"path/filepath"
	"strings"
)

var App IApp
var DefaultHomeDir string
var defaultP2pPort = "20000"
var DefaultRpcPort = "20001"
var defaultPrivateFile = "key.json"
var defaultKey = "future_mine_chain"
var defaultExternalIp = "0.0.0.0"
var DefaultFallBack = int64(-1)

// Config is the node startup parameter
type Config struct {
	ConfigFile string `long:"config" description:"Start with a configuration file"`
	Home       string `long:"appdata" description:"Path to application home directory"`
	Data       string `long:"data" description:"Path to application data directory"`
	Logging    bool   `long:"logging" description:"Logging switch"`
	ExternalIp string `long:"externalip" description:"External network IP address"`
	Bootstrap  string `long:"bootstrap" description:"Custom bootstrap"`
	P2PPort    string `long:"p2pport" description:"Add an interface/port to listen for connections"`
	RpcPort    string `long:"rpcport" description:"Add an interface/port to listen for RPC connections"`
	RpcTLS     bool   `long:"rpctls" description:"Open TLS for the RPC server -- NOTE: This is only allowed if the RPC server is bound to localhost"`
	RpcCert    string `long:"rpccert" description:"File containing the certificate file"`
	RpcKey     string `long:"rpckey" description:"File containing the certificate key"`
	RpcPass    string `long:"rpcpass" description:"Password for RPC connections"`
	TestNet    bool   `long:"testnet" description:"Use the test network"`
	KeyFile    string `long:"keyfile" description:"If you participate in mining, you need to configure the mining address key file"`
	KeyPass    string `long:"keypass" description:"The decryption password for key file"`
	FallBackTo int64  `long:"fallbackto" description:"Force back to a height"`
}

// LoadConfig load the parse node startup parameter
func LoadConfig(app IApp) (*Config, error) {
	App = app
	DefaultHomeDir = utils.AppDataDir(App.AppName(), false)
	cfg := &Config{
		Home:       DefaultHomeDir,
		P2PPort:    defaultP2pPort,
		RpcPort:    DefaultRpcPort,
		FallBackTo: DefaultFallBack,
	}
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	preParser := newConfigParser(cfg, flags.HelpFlag)
	_, err := preParser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type != flags.ErrHelp {
			return nil, err
		} else if ok && e.Type == flags.ErrHelp {
			return nil, err
		}
	}

	if cfg.ConfigFile != "" {
		_, err = toml.DecodeFile(cfg.ConfigFile, cfg)
		if err != nil {
			return nil, err
		}
	}

	// Set the default external IP. If the external IP is not set,
	// other nodes can only know you but cannot send messages to you.
	if cfg.ExternalIp == "" {
		cfg.ExternalIp = defaultExternalIp
	}

	// Node data and file storage directory, if not set,
	// use the default directory
	if cfg.Home == "" {
		cfg.Home = DefaultHomeDir
	}

	// p2p service listening port, if not, use the default port
	if cfg.P2PPort == "" {
		cfg.P2PPort = defaultP2pPort
	}

	// rpc service listening port, if not, use the default port
	if cfg.RpcPort == "" {
		cfg.RpcPort = DefaultRpcPort
	}

	if cfg.TestNet {
		App.InitTestNet()
	}

	// p2p same network label, the label is different and cannot communicate
	app.InitP2pNet()

	if !utils.IsExist(cfg.Home) {
		if err := os.Mkdir(cfg.Home, os.ModePerm); err != nil {
			return nil, err
		}
	}
	if cfg.Data == "" {
		cfg.Data = cfg.Home + "/" + cfg.P2PPort
	}
	if !utils.IsExist(cfg.Data) {
		if err := os.Mkdir(cfg.Data, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// Each node requires a secp256k1 private key, which is used as the p2p id
	// generation and signature of the node that generates the block.
	// If this parameter is not configured in the startup parameter,
	// the node will be automatically generated and loaded automatically at startup
	/*if cfg.KeyFile == "" {
		cfg.KeyFile = defaultPrivateFile
		cfg.NodePrivate, err = LoadNodePrivate(cfg.DataDir+"/"+cfg.KeyFile, defaultKey)
		if err != nil {
			cfg.NodePrivate, err = CreateNewNodePrivate(app.NetWork())
			if err != nil {
				return nil, fmt.Errorf("create new node priavte failed! %s", err.Error())
			}
		}
		j, err := keystore.PrivateToJson(App.NetWork(), cfg.NodePrivate.PrivateKey, cfg.NodePrivate.Mnemonic, []byte(defaultKey))
		if err != nil {
			return nil, fmt.Errorf("key json creation failed! %s", err.Error())
		}
		bytes, _ := json.Marshal(j)
		err = ioutil.WriteFile(cfg.DataDir+"/"+cfg.KeyFile, bytes, 0644)
		if err != nil {
			return nil, fmt.Errorf("write jsonfile failed! %s", err.Error())
		}
	} else {
		// The private key of the node is encrypted in the key file,
		// and a password is required to unlock the key file
		if cfg.KeyPass == "" {
			fmt.Println("Please enter the password for the keyfile:")
			passWd, err := readPassWd()
			if err != nil {
				return nil, fmt.Errorf("read password failed! %s", err.Error())
			}
			cfg.KeyPass = string(passWd)
		}
		cfg.NodePrivate, err = LoadNodePrivate(cfg.KeyFile, cfg.KeyPass)
		if err != nil {
			return nil, fmt.Errorf("failed to load keyfile %s! %s", cfg.KeyFile, err.Error())
		}
	}*/

	// If this parameter is true, the log is also written to the file
	if cfg.Logging {
		logDir := cfg.Data + "/log"
		if !utils.IsExist(logDir) {
			if err := os.Mkdir(logDir, os.ModePerm); err != nil {
				return nil, err
			}
		}
		utils.CleanAndExpandPath(logDir)
		logDir = filepath.Join(logDir, App.NetWork())
		log2.InitLogRotator(filepath.Join(logDir, "future_mine.log"))
	}
	log.Info("Data storage directory", "module", "config", "path", cfg.Data)
	return cfg, nil
}

func newConfigParser(cfg *Config, options flags.Options) *flags.Parser {
	parser := flags.NewParser(cfg, options)
	return parser
}
