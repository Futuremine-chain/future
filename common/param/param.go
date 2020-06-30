package param

import (
	"github.com/Futuremine-chain/futuremine/common/private"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"time"
)

const (
	// Block interval period
	BlockInterval = uint64(5)
	// Re-election interval
	CycleInterval = 60 * 60 * 24
	// Maximum number of super nodes
	SuperSize = 3

	DPosSize = SuperSize*2/3 + 1
)

const (
	// Mainnet logo
	MainNet = "mainnet"
	// Testnet logo
	TestNet = "testnet"

	Version = "0.1.1"
)

const (
	MaxReadBytes = 1024 * 10
	MaxReqBytes  = MaxReadBytes * 1000
)

type Param struct {
	Name              string
	Data              string
	App               string
	FallBack          int64
	PubKeyHashAddrID  [2]byte
	PubKeyHashTokenID [2]byte
	Logging           bool
	PeerRequestChan   uint32
	*PrivateParam
	*TokenParam
	*P2pParam
	*RpcParam
	*DPosParam
	*PoolParam
	private.IPrivate
}

type TokenParam struct {
	CoinBase     uint64
	Proportion   uint64
	MinCoinCount uint64
	MaxCoinCount uint64
	MainToken    arry.Address
}

type PrivateParam struct {
	PrivateFile string
	PrivatePass string
}

type P2pParam struct {
	P2pPort    string
	ExternalIp string
	NetWork    string
	CustomBoot string
}

type RpcParam struct {
	RpcIp      string
	RpcPort    string
	RpcTLS     bool
	RpcCert    string
	RpcCertKey string
	RpcPass    string
}

type DPosParam struct {
	BlockInterval uint64
	CycleInterval uint64
	SuperSize     int
	DPosSize      int
	GenesisTime   uint64
	GenesisCycle  uint64
}

type PoolParam struct {
	MsgExpiredTime     int64
	MonitorMsgInterval time.Duration
	MaxPoolMsg         int
	MaxAddressMsg      uint64
}

var TestNetParam = &Param{
	Name:              TestNet,
	Data:              "data",
	App:               "future-mine-chain",
	FallBack:          -1,
	PubKeyHashAddrID:  [2]byte{0x1f, 0x13},
	PubKeyHashTokenID: [2]byte{0x0f, 0x03},
	Logging:           true,
	PeerRequestChan:   1000,
	PrivateParam: &PrivateParam{
		PrivateFile: "key.json",
		PrivatePass: "fmc",
	},
	TokenParam: &TokenParam{
		CoinBase:     10 * 1e8,
		Proportion:   10000,
		MinCoinCount: 1 * 1e4,
		MaxCoinCount: 1 * 1e10 * 1e8,
		MainToken:    arry.StringToAddress("FMC"),
	},
	P2pParam: &P2pParam{
		NetWork:    TestNet + "_FUTURE_MINE_CHAIN",
		P2pPort:    "19160",
		ExternalIp: "0.0.0.0",
		CustomBoot: "",
	},
	RpcParam: &RpcParam{
		RpcIp:      "127.0.0.1",
		RpcPort:    "19161",
		RpcTLS:     false,
		RpcCert:    "",
		RpcCertKey: "",
		RpcPass:    "",
	},
	DPosParam: &DPosParam{
		BlockInterval: BlockInterval,
		CycleInterval: CycleInterval,
		SuperSize:     CycleInterval,
		DPosSize:      DPosSize,
		GenesisTime:   1592268410,
		GenesisCycle:  1592268410 / CycleInterval,
	},
	PoolParam: &PoolParam{
		MaxPoolMsg:         100000,
		MsgExpiredTime:     60 * 60 * 3,
		MonitorMsgInterval: 2,
		MaxAddressMsg:      1000,
	},
}

var MainNetParam = &Param{
	Name:              MainNet,
	Data:              "data",
	App:               "future-mine-chain",
	FallBack:          -1,
	PubKeyHashAddrID:  [2]byte{0x1e, 0x12},
	PubKeyHashTokenID: [2]byte{0x0e, 0x02},
	Logging:           true,
	PeerRequestChan:   1000,
	PrivateParam: &PrivateParam{
		PrivateFile: "key.json",
		PrivatePass: "fmc",
	},
	TokenParam: &TokenParam{
		CoinBase:     10 * 1e8,
		Proportion:   10000,
		MinCoinCount: 1 * 1e4,
		MaxCoinCount: 1 * 1e10 * 1e8,
		MainToken:    arry.StringToAddress("FMC"),
	},
	P2pParam: &P2pParam{
		NetWork:    MainNet + "_FUTURE_MINE_CHAIN",
		P2pPort:    "29160",
		ExternalIp: "0.0.0.0",
		CustomBoot: "",
	},
	RpcParam: &RpcParam{
		RpcIp:      "127.0.0.1",
		RpcPort:    "29161",
		RpcTLS:     false,
		RpcCert:    "",
		RpcCertKey: "",
		RpcPass:    "",
	},
	DPosParam: &DPosParam{
		BlockInterval: BlockInterval,
		CycleInterval: CycleInterval,
		SuperSize:     CycleInterval,
		DPosSize:      DPosSize,
		GenesisTime:   1592268410,
		GenesisCycle:  1592268410 / CycleInterval,
	},
	PoolParam: &PoolParam{
		MaxPoolMsg:         100000,
		MsgExpiredTime:     60 * 60 * 3,
		MonitorMsgInterval: 2,
		MaxAddressMsg:      1000,
	},
}
