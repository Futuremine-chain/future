package server

type IService interface {
	Name() string
	Start() error
	Stop() error
}

type IServer interface {
	Register(IService)
	Start() error
	Stop() error
}
