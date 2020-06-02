package config

type IApp interface {
	AppName() string
	Version() string
	NetWork() string
	P2pNetWork() string
	TestNet() string
	MainNet() string
	Setting() *Config
	InitTestNet()
	InitP2pNet()
	InitSetting(*Config)
}
