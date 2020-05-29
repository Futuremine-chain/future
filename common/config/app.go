package config

type IApp interface {
	AppName() string
	Version() string
	NetWork() string
	P2pNetWork() string
	TestNet()string
	MainNet()string
	InitTestNet()
	InitP2pNet()
}
