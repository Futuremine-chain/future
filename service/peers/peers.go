package peers

import (
	log "github.com/Futuremine-chain/futuremine/tools/log/log15"
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
	local  *Peer
	cache  map[string]*Peer
	idList []string
	rwm    sync.RWMutex
}

func NewPeers() *Peers {
	return &Peers{cache: make(map[string]*Peer, maxPeers)}
}

func (p *Peers) Name() string {
	return module
}

func (p *Peers) Start() error {
	log.Info("Peers started successfully", "module", module)
	go p.Monitoring()
	return nil
}

func (p *Peers) Stop() error {
	return nil
}

func (p *Peers) AddressExist(address *peer.AddrInfo) bool {
	p.rwm.RLock()
	defer p.rwm.RUnlock()

	if peer, ok := p.cache[address.ID.String()]; !ok {
		return false
	} else {
		if peer.Address.String() != address.String() {
			return false
		}
	}
	return true
}

func (p *Peers) AddPeer(peer *Peer) {
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

func (p *Peers) Monitoring() {
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

func (p *Peers) isAlive(peer *Peer) bool {
	stream, err := peer.Conn.Create(peer.Address.ID)
	if err != nil {
		return false
	}
	stream.Reset()
	stream.Close()
	return true
}

func (p *Peers) RandomPeer() *Peer {
	p.rwm.Lock()
	defer p.rwm.Unlock()

	if len(p.idList) == 0 {
		return nil
	}
	index := rand.New(rand.NewSource(time.Now().Unix())).Int31n(int32(len(p.idList)))
	peerId := p.idList[index]
	return p.cache[peerId]
}

func (p *Peers) Local() *Peer {
	return p.local
}

func (p *Peers) SetLocal(local *Peer) {
	p.local = local
}

func (p *Peers) PeersMap() map[string]*Peer {
	p.rwm.RLock()
	defer p.rwm.RUnlock()

	re := make(map[string]*Peer)
	for key, value := range p.cache {
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
