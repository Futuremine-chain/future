package request

import (
	"errors"
	"github.com/Futuremine-chain/futuremine/server"
	"github.com/Futuremine-chain/futuremine/service/peers"
	"github.com/Futuremine-chain/futuremine/types"
	"github.com/libp2p/go-libp2p-core/network"
)

var (
	Err_BlockNotFound = errors.New("block not exist")
	Err_PeerClosed    = errors.New("peer has closed")
)

type IRequestHandler interface {
	server.IService
	ISend
	IRegister
	IResponse
}

type ISend interface {
	LastHeight(conn *peers.Conn) (uint64, error)
	SendMsg(conn *peers.Conn, msg types.IMessage) error
	SendBlock(conn *peers.Conn, block types.IBlock) error
	GetBlocks(conn *peers.Conn, height uint64) ([]types.IBlock, error)
	GetBlock(conn *peers.Conn, height uint64) (types.IBlock, error)
	IsEqual(conn *peers.Conn, header types.IHeader) (bool, error)
}

type IRegister interface {
	RegisterReceiveBlock(func(types.IBlock) error)
	RegisterReceiveMessage(func(types.IMessage) error)
}

type IResponse interface {
	SendToReady(stream network.Stream)
}
