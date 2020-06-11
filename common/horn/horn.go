package horn

import (
	"github.com/Futuremine-chain/futuremine/service/gorutinue"
	"github.com/Futuremine-chain/futuremine/service/peers"
	"github.com/Futuremine-chain/futuremine/service/request"
	log "github.com/Futuremine-chain/futuremine/tools/log/log15"
	"github.com/Futuremine-chain/futuremine/types"
)

const module = "horn"

type Horn struct {
	peers   *peers.Peers
	request request.IRequestHandler
	local   *peers.Peer
	gPool   *gorutinue.Pool
}

func NewHorn(peers *peers.Peers, gPool *gorutinue.Pool) *Horn {
	return &Horn{
		peers:   peers,
		request: nil,
		local:   peers.Local(),
		gPool:   gPool,
	}
}

func (h *Horn) BroadcastMsg(message types.IMessage) {
	peers := h.peers.PeersMap()
	for id, peer := range peers {
		if id != h.local.Address.ID.String() {
			if err := h.gPool.AddTask(gorutinue.NewTask(
				func() error {
					return h.request.SendMsg(peer.Conn, message)
				})); err != nil {
				log.Warn("Adding the task to send the message failed", "module", module,
					"hash", message.Hash().String(), "target", peer.Address.String())
			}
		}
	}
}

func (h *Horn) BroadcastBlock(block types.IBlock) {
	peers := h.peers.PeersMap()
	for id, peer := range peers {
		if id != h.local.Address.ID.String() {
			if err := h.gPool.AddTask(gorutinue.NewTask(
				func() error {
					return h.request.SendBlock(peer.Conn, block)
				})); err != nil {
				log.Warn("Adding the task to send the block failed", "module", module,
					"height", block.Hash().String(), "target", peer.Address.String())
			}
		}
	}
}
