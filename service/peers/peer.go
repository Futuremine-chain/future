package peers

import (
	"crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
)

type Peer struct {
	Private crypto.PrivateKey
	Address *peer.AddrInfo
	Conn    *Conn
}

func NewPeer(private crypto.PrivateKey, addr *peer.AddrInfo, createF CreateConnF) *Peer {
	return &Peer{Private: private, Address: addr, Conn: &Conn{PeerId: addr.ID, Create: createF}}
}

type CreateConnF func(peerId peer.ID) (network.Stream, error)

type Conn struct {
	Stream network.Stream
	PeerId peer.ID
	Create CreateConnF
}

func (s *Conn) Close() {
	s.Stream.Reset()
	s.Stream.Close()
	s.Stream = nil
}
