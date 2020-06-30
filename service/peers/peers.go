package peers

import (
	request2 "github.com/Futuremine-chain/futuremine/service/request"
	log "github.com/Futuremine-chain/futuremine/tools/log/log15"
	"github.com/Futuremine-chain/futuremine/types"
	"github.com/libp2p/go-libp2p-core/peer"
	"math/rand"
	"sync"
	"time"
)

const (
	module             = "peers"
	maxPeers           = 1000000
	monitoringInterval = 10
)

type Peers struct {
	local      *types.Peer
	cache      map[string]*types.Peer
	idList     []string
	rwm        sync.RWMutex
	close      chan bool
	peerInfo   map[string]*types.Local
	infoWm     sync.RWMutex
	reqHandler request2.IRequestHandler
}

func NewPeers(reqHandler request2.IRequestHandler) *Peers {
	return &Peers{
		cache:      make(map[string]*types.Peer, maxPeers),
		close:      make(chan bool),
		peerInfo:   make(map[string]*types.Local),
		reqHandler: reqHandler,
	}
}

func (p *Peers) Name() string {
	return module
}

func (p *Peers) Start() error {
	log.Info("Peers started successfully", "module", module)
	go p.monitoring()
	go p.getPeerLocal()
	return nil
}

func (p *Peers) Stop() error {
	log.Info("Peers was stopped", "module", module)
	return nil
}

func (p *Peers) Info() map[string]interface{} {
	return map[string]interface{}{
		"connections": p.Count(),
	}
}

func (p *Peers) AddressExist(address *peer.AddrInfo) bool {
	p.rwm.RLock()
	defer p.rwm.RUnlock()

	if _, ok := p.cache[address.ID.String()]; !ok {
		return false
	}
	return true
}

func (p *Peers) AddPeer(peer *types.Peer) {
	p.rwm.Lock()
	defer p.rwm.Unlock()

	if len(p.cache) >= maxPeers {
		return
	}
	p.cache[peer.Address.ID.String()] = peer
	p.idList = append(p.idList, peer.Address.ID.String())
	log.Info("Add a peer", "module", module, "id", peer.Address.ID.String(), "address", peer.Address.String())
}

func (p *Peers) RemovePeer(reId string) {
	p.rwm.Lock()
	defer p.rwm.Unlock()

	for index, id := range p.idList {
		if id == reId {
			p.idList = append(p.idList[0:index], p.idList[index+1:]...)
			delete(p.cache, reId)
			log.Info("Delete a peer", "id", reId)
			break
		}
	}
}

func (p *Peers) monitoring() {
	t := time.NewTicker(time.Second * monitoringInterval)
	defer t.Stop()

	for range t.C {
		for id, peer := range p.cache {
			if id != p.local.Address.ID.String() {
				if !p.isAlive(peer) {
					p.RemovePeer(id)
				}
			}
		}
	}
}

func (p *Peers) getPeerLocal() {
	t := time.NewTicker(time.Second * 120)
	defer t.Stop()

	for range t.C {
		peers := p.PeersMap()
		for id, peer := range peers {
			if id != p.local.Address.ID.String() {
				local, err := p.reqHandler.LocalInfo(peer.Conn)
				if err != nil {
					p.infoWm.Lock()
					p.peerInfo[id] = local
					p.infoWm.Unlock()
				}
			}
		}
	}
}

func (p *Peers) isAlive(peer *types.Peer) bool {
	stream, err := peer.Conn.Create(peer.Address.ID)
	if err != nil {
		return false
	}
	stream.Reset()
	stream.Close()
	return true
}

func (p *Peers) RandomPeer() *types.Peer {
	p.rwm.Lock()
	defer p.rwm.Unlock()

	if len(p.idList) == 0 {
		return nil
	}
	index := rand.New(rand.NewSource(time.Now().Unix())).Int31n(int32(len(p.idList)))
	peerId := p.idList[index]
	return p.cache[peerId]
}

func (p *Peers) Local() *types.Peer {
	return p.local
}

func (p *Peers) SetLocal(local *types.Peer) {
	p.local = local
}

func (p *Peers) PeersMap() map[string]*types.Peer {
	p.rwm.RLock()
	defer p.rwm.RUnlock()

	re := make(map[string]*types.Peer)
	for key, value := range p.cache {
		re[key] = value
	}
	return re
}

func (p *Peers) PeersInfo() map[string]*types.Local {
	p.infoWm.RLock()
	defer p.infoWm.RUnlock()

	re := make(map[string]*types.Local)
	for key, value := range p.peerInfo {
		re[key] = value
	}
	return re
}

func (p *Peers) Count() uint32 {
	p.rwm.RLock()
	defer p.rwm.RUnlock()

	count := uint32(len(p.cache))
	return count
}

func (p *Peers) Peer(id string) *types.Peer {
	p.rwm.RLock()
	defer p.rwm.RUnlock()

	peer := p.cache[id]
	return peer
}
