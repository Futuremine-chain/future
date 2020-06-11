package config

import "github.com/Futuremine-chain/futuremine/tools/arry"

type IApp interface {
	AppName() string
	Version() string
	NetWork() string
	P2pNetWork() string
	TestNet() string
	MainNet() string
	MainToken() arry.Address
	Setting() *Config
	InitTestNet()
	InitP2pNet()
	InitSetting(*Config)
}
