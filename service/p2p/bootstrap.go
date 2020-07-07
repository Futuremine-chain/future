package p2p

import (
	"context"
	"github.com/Futuremine-chain/futuremine/tools/crypto/ecc/secp256k1"
	log "github.com/Futuremine-chain/futuremine/tools/log/log15"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
	"strings"
)

// Default boot node list
var DefaultBootPeers []multiaddr.Multiaddr

// Custom boot node list
var CustomBootPeers []multiaddr.Multiaddr

func init() {
	for _, s := range []string{
		"/ip4/117.34.92.9/tcp/19100/ipfs/16Uiu2HAmKxgLc1zsbbGP4CobP2mvjeqzLzBgXQNtHcHcRese6Jtv",
	} {
		ma, err := multiaddr.NewMultiaddr(s)
		if err != nil {
			panic(err)
		}
		DefaultBootPeers = append(DefaultBootPeers, ma)
	}
}

func IsBootPeers(id peer.ID) bool {
	bootstrap := DefaultBootPeers
	if len(CustomBootPeers) > 0 {
		bootstrap = CustomBootPeers
	}
	for _, bootstrap := range bootstrap {
		if id.String() == strings.Split(bootstrap.String(), "/")[6] {
			return true
		}
	}
	return false
}

func NewBoot(ip, port, external string, private *secp256k1.PrivateKey) (*P2p, error) {
	host, err := NewP2PHost(private, ip, port, external)
	if err != nil {
		return nil, err
	}
	p2p := &P2p{host: host}
	log.Info("Host created", "id", p2p.host.ID(), "address", p2p.host.Addrs())
	return p2p, nil
}

func (p *P2p) StartBoot() error {
	var err error
	p.dht, err = dht.New(context.Background(), p.host)
	if err != nil {
		return err
	}
	log.Info("Start the boot node", "module", module)
	if err = p.dht.Bootstrap(context.Background()); err != nil {
		return err
	}
	return nil
}
