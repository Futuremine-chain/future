package p2p

import (
	"context"
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
		"/ip4/127.0.0.1/tcp/2211/ipfs/16Uiu2HAmJBN2NPzA1nHpyN2Doifq5CK36GqRf2vwcQ1gcz4JoZRm",
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
