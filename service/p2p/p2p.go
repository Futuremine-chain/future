package p2p

import (
	"context"
	"crypto"
	"fmt"
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/common/private"
	"github.com/Futuremine-chain/futuremine/service/peers"
	"github.com/Futuremine-chain/futuremine/service/request"
	"github.com/Futuremine-chain/futuremine/tools/crypto/ecc/secp256k1"
	log "github.com/Futuremine-chain/futuremine/tools/log/log15"
	"github.com/Futuremine-chain/futuremine/tools/utils"
	"github.com/libp2p/go-libp2p"
	core "github.com/libp2p/go-libp2p-core"
	crypto2 "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	discovery "github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	p2pcfg "github.com/libp2p/go-libp2p/config"
	"github.com/multiformats/go-multiaddr"
	"sync"
	"time"
)

const module = "p2p"

type P2p struct {
	host       core.Host
	local      *peers.Peer
	dht        *dht.IpfsDHT
	peers      *peers.Peers
	reqHandler request.IRequestHandler
	close      chan bool
	closed     chan bool
}

func NewP2p(cfg *config.Config, ps *peers.Peers, reqHandler request.IRequestHandler, priv *secp256k1.PrivateKey) (*P2p, error) {
	var err error
	ser := &P2p{
		peers:      ps,
		reqHandler: reqHandler,
		close:      make(chan bool),
		closed:     make(chan bool),
	}
	if cfg.Boot != "" {
		ma, err := multiaddr.NewMultiaddr(cfg.Boot)
		if err != nil {
			return nil, fmt.Errorf("incorrect bootstrap node address， %s", err)
		}
		CustomBootPeers = append(CustomBootPeers, ma)
	}

	host, err := newP2PHost(priv, cfg.ExternalIp, cfg.P2PPort, cfg.ExternalIp)
	if err != nil {
		return nil, err
	}
	ser.host = host
	ser.local = peers.NewPeer(priv,
		&peer.AddrInfo{
			ID:    host.ID(),
			Addrs: host.Addrs()}, nil)
	ser.initP2pHandle()
	ps.SetLocal(ser.local)
	log.Info("P2p host created", "module", module, "id", host.ID(), "address", host.Addrs())
	return ser, nil
}

func newP2PHost(private *secp256k1.PrivateKey, ip, port, external string) (core.Host, error) {
	ips := utils.GetLocalIp()
	ips = append(ips, external)
	f := newFactory(ips, port)
	p2pKey, err := crypto2.UnmarshalSecp256k1PrivateKey(private.Serialize())
	if err != nil {
		return nil, err
	}
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", port)),
		libp2p.Identity(p2pKey),
		libp2p.DefaultMuxers,
		libp2p.EnableRelay(),
		libp2p.AddrsFactory(f),
	}
	return libp2p.New(context.Background(), opts...)
}

func newFactory(ips []string, port string) p2pcfg.AddrsFactory {
	return func(addrs []multiaddr.Multiaddr) []multiaddr.Multiaddr {
		addrs = []multiaddr.Multiaddr{}
		for _, ip := range ips {
			extMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%s", ip, port))
			if extMultiAddr != nil {
				addrs = append(addrs, extMultiAddr)
			}
		}
		return addrs
	}
}

func (p *P2p) Name() string {
	return module
}

func (p *P2p) Start() error {
	if err := p.connectBootNode(); err != nil {
		log.Error("Failed to connect the boot node!", "module", module, "error", err)
		return err
	}

	go p.peerDiscovery()
	log.Info("P2P started successfully", "module", module)
	return nil
}

func (p *P2p) Stop() error {
	p.close <- true
	<-p.closed
	if err := p.host.Close(); err != nil {
		return err
	}
	log.Info("Stop P2P server")
	return nil
}

func (p *P2p) Addr() string {
	addrs := p.host.Addrs()
	var rs string
	for _, addr := range addrs {
		rs += "[" + addr.String() + "]"
	}
	return rs
}

func (p *P2p) ID() string {
	return p.host.ID().String()
}

func (p *P2p) newStream(id peer.ID) (network.Stream, error) {
	return p.host.NewStream(context.Background(), id, protocol.ID(config.App.P2pNetWork()))
}

func (p *P2p) initP2pHandle() {
	p.host.SetStreamHandler(protocol.ID(config.App.P2pNetWork()), p.reqHandler.SendToReady)
}

func (p *P2p) connectBootNode() error {
	var err error
	p.dht, err = dht.New(context.Background(), p.host)
	if err != nil {
		return err
	}

	log.Info("Initializing node DHT", "module", module)
	if err = p.dht.Bootstrap(context.Background()); err != nil {
		return err
	}

	boots := DefaultBootPeers
	if len(CustomBootPeers) > 0 {
		boots = CustomBootPeers
	}
	var wg sync.WaitGroup
	for _, address := range boots {
		addrInfo, _ := peer.AddrInfoFromP2pAddr(address)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := p.host.Connect(context.Background(), *addrInfo); err != nil {
				log.Warn("Failed to connection established with boot node", "module", module, "error", err)
			} else {
				log.Info("Connection established with boot node", "module", module, "peer", *addrInfo)
			}
		}()
	}
	wg.Wait()
	return nil
}

// peerDiscovery new nodes every 8s
func (p *P2p) peerDiscovery() {
	rouDis := discovery.NewRoutingDiscovery(p.dht)
	discovery.Advertise(context.Background(), rouDis, config.App.P2pNetWork())

	for {
		select {
		case _ = <-p.close:
			p.closed <- true
			return
		default:
			log.Info("Look for other peers...", "module", module)
			ch, err := rouDis.FindPeers(context.Background(), config.App.P2pNetWork())
			if err != nil {
				log.Error("Peer search failed", "module", module, "error", err)
				time.Sleep(time.Second * 10)
				continue
			}
			p.readAddrInfo(ch)
		}
		time.Sleep(time.Second * 8)
	}
}

func (p *P2p) readAddrInfo(addrCh <-chan peer.AddrInfo) {
	for {
		select {
		case addrInfo, ok := <-addrCh:
			if ok {
				if addrInfo.ID == p.local.Address.ID || IsBootPeers(addrInfo.ID) {
					continue
				}
				if !p.peers.AddressExist(&addrInfo) {
					if !p.isAlive(addrInfo.ID) {
						p.peers.RemovePeer(addrInfo.ID.String())
						continue
					}
					p.peers.AddPeer(peers.NewPeer(nil, cpAddrInfo(&addrInfo), p.newStream))
				}
			} else {
				return
			}
		}
	}
}

func (p *P2p) isAlive(id peer.ID) bool {
	stream, err := p.newStream(id)
	if err != nil {
		return false
	}
	stream.Reset()
	stream.Close()
	return true
}

func PrivateToP2pId(key private.Private) (peer.ID, error) {
	p2pPriKey, err := crypto.UnmarshalSecp256k1PrivateKey(key.Serialize())
	if err != nil {
		return "", err
	}
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", "65535")),
		libp2p.Identity(p2pPriKey),
	}
	host, err := libp2p.New(context.Background(), opts...)
	if err != nil {
		return "", err
	}
	defer host.Close()
	return host.ID(), nil
}

func cpAddrInfo(addr *peer.AddrInfo) *peer.AddrInfo {
	bytes, _ := addr.MarshalJSON()
	destAddr := new(peer.AddrInfo)
	destAddr.UnmarshalJSON(bytes)
	return destAddr
}
