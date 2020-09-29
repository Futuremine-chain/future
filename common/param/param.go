package param

import (
	"github.com/Futuremine-chain/future/common/private"
	"github.com/Futuremine-chain/future/tools/arry"
	"time"
)

const (
	// Block interval period
	BlockInterval = uint64(30)
	// Re-election interval
	CycleInterval = 60 * 60 * 24 * 3
	// Maximum number of super nodes
	SuperSize = 9

	DPosSize = SuperSize*2/3 + 1
)

const (
	// Mainnet logo
	MainNet = "mainnet"
	// Testnet logo
	TestNet = "testnet"

	Version = "0.2.6"
)

const (
	MaxReadBytes = 1024 * 10
	MaxReqBytes  = MaxReadBytes * 1000
)

// AtomsPerCoin is the number of atomic units in one coin.
const AtomsPerCoin = 1e8

type Param struct {
	Name              string
	Data              string
	App               string
	RollBack          uint64
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
	PreCirculation        uint64
	Circulation           uint64
	CoinBase              float64
	CoinCoefficient       float64
	EveryChangeCoinHeight uint64
	Proportion            uint64
	MinCoinCount          float64
	MaxCoinCount          float64
	MinimumTransfer       uint64
	Consume               uint64
	MainToken             arry.Address
	EaterAddress          arry.Address
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
	App:               "future-chain",
	RollBack:          0,
	PubKeyHashAddrID:  [2]byte{0x08, 0x51},
	PubKeyHashTokenID: [2]byte{0x08, 0x62},
	Logging:           true,
	PeerRequestChan:   1000,
	PrivateParam: &PrivateParam{
		PrivateFile: "key.json",
		PrivatePass: "fm",
	},
	TokenParam: &TokenParam{
		PreCirculation:        38000000 * AtomsPerCoin,
		Circulation:           248034687.5 * AtomsPerCoin,
		CoinBase:              100 * AtomsPerCoin,
		Consume:               1e4 * AtomsPerCoin,
		EveryChangeCoinHeight: 1051201,
		CoinCoefficient:       -0.5,
		Proportion:            10000,
		MinCoinCount:          1 * 1e4,
		MaxCoinCount:          9 * 1e10,
		MinimumTransfer:       0.0001 * AtomsPerCoin,
		MainToken:             arry.StringToAddress("FM"),
		EaterAddress:          arry.StringToAddress("FmcoinEaterAddressDontSend000000000"),
	},
	P2pParam: &P2pParam{
		NetWork:    TestNet + "_FUTURE_CHAIN",
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
		SuperSize:     SuperSize,
		DPosSize:      DPosSize,
		GenesisTime:   1592268410,
		GenesisCycle:  1592268410 / CycleInterval,
	},
	PoolParam: &PoolParam{
		MaxPoolMsg:         100000,
		MsgExpiredTime:     60 * 60 * 3,
		MonitorMsgInterval: 10,
		MaxAddressMsg:      1000,
	},
}

var MainNetParam = &Param{
	Name:              MainNet,
	Data:              "data",
	App:               "future-chain",
	RollBack:          0,
	PubKeyHashAddrID:  [2]byte{0x08, 0x15},
	PubKeyHashTokenID: [2]byte{0x08, 0x24},
	Logging:           true,
	PeerRequestChan:   1000,
	PrivateParam: &PrivateParam{
		PrivateFile: "key.json",
		PrivatePass: "fm",
	},
	TokenParam: &TokenParam{
		PreCirculation:        38000000 * AtomsPerCoin,
		Circulation:           248034687.5 * AtomsPerCoin,
		CoinBase:              100 * AtomsPerCoin,
		Consume:               1e4 * AtomsPerCoin,
		EveryChangeCoinHeight: 1051201,
		CoinCoefficient:       -0.5,
		Proportion:            10000,
		MinCoinCount:          1 * 1e4,
		MaxCoinCount:          9 * 1e10,
		MinimumTransfer:       0.0001 * AtomsPerCoin,
		MainToken:             arry.StringToAddress("FM"),
		EaterAddress:          arry.StringToAddress("FMcoinEaterAddressDontSend000000000"),
	},
	P2pParam: &P2pParam{
		NetWork:    MainNet + "_FUTURE_CHAIN",
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
		SuperSize:     SuperSize,
		DPosSize:      DPosSize,
		GenesisTime:   1592268410,
		GenesisCycle:  1592268410 / CycleInterval,
	},
	PoolParam: &PoolParam{
		MaxPoolMsg:         100000,
		MsgExpiredTime:     60 * 60 * 3,
		MonitorMsgInterval: 10,
		MaxAddressMsg:      1000,
	},
}

type PreCirculation struct {
	Address string
	Note    string
	Amount  uint64
}

var PreCirculations = []PreCirculation{
	{
		Address: "FMhJy6XEUKR2hYJ8SdWkoi5YosVxzjqdEyu",
		Note:    "",
		Amount:  38000000 * AtomsPerCoin,
	},
}
