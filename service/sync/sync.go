package sync

import (
	"errors"
	"github.com/Futuremine-chain/futuremine/common/blockchain"
	"github.com/Futuremine-chain/futuremine/common/dpos"
	"github.com/Futuremine-chain/futuremine/service/peers"
	"github.com/Futuremine-chain/futuremine/service/request"
	log "github.com/Futuremine-chain/futuremine/tools/log/log15"
	"github.com/Futuremine-chain/futuremine/types"
	"time"
)

const module = "sync"

type Sync struct {
	chain   blockchain.IChain
	request request.IRequestHandler
	peers   *peers.Peers
	curPeer *peers.Peer
	dPos    dpos.IDPos
	stop    chan bool
	stopped chan bool
}

func NewSync(peers *peers.Peers, dPos dpos.IDPos, request request.IRequestHandler, chain blockchain.IChain) *Sync {
	s := &Sync{
		chain:   chain,
		peers:   peers,
		dPos:    dPos,
		request: request,
		stop:    make(chan bool),
		stopped: make(chan bool),
	}
	return s
}

func (s *Sync) Name() string {
	return module
}

func (s *Sync) Start() error {
	//go s.syncBlocks()
	log.Info("Sync started successfully", "module", module)
	return nil
}

func (s *Sync) Stop() error {
	close(s.stop)
	<-s.stopped
	log.Info("Stop sync block", "module", module)
	return nil
}

// Start sync block
func (s *Sync) syncBlocks() {
	for {
		select {
		case _, _ = <-s.stop:
			s.stopped <- true
			return
		default:
			s.createSyncStream()
			s.syncFromConn()

		}
		time.Sleep(time.Millisecond * 1000)
	}
}

// Create a network channel of the synchronization block, and randomly
// select a new peer node for synchronization every 1s.
func (s *Sync) createSyncStream() {
	for {
		select {
		case _, _ = <-s.stop:
			return
		default:
			s.findSyncPeer()
			return
		}
	}
}

// Replace the new peer node
func (s *Sync) findSyncPeer() {
	t := time.NewTicker(time.Second * 10)
	defer t.Stop()

	for {
		select {
		case _, _ = <-s.stop:
			return
		case _ = <-t.C:
			s.curPeer = s.peers.RandomPeer()
			if s.curPeer == nil {
				//log.Warn("No available peers were found， wait...")
			} else {
				//log.Info("Find an available peer", "peer", bm.syncPeer)
				return
			}
		}
	}
}

// Synchronize blocks from the stream and verify storage
func (s *Sync) syncFromConn() error {
	for {
		select {
		case _, _ = <-s.stop:
			return nil
		default:
			if s.curPeer == nil {
				return errors.New("no current peer")
			}
			localHeight := s.chain.LastHeight()

			// Get the block of the remote node from the next block height，
			// If the error is that the peer has stopped, delete the peer.
			// If the storage fails locally, the remote block verification
			// is performed, the verification proves that the local block
			// is wrong, and the local chain is rolled back to the valid block.
			blocks, err := s.request.GetBlocks(s.curPeer.Conn, localHeight+1)
			if err != nil {
				if err == request.Err_PeerClosed {
					s.peers.RemovePeer(s.curPeer.Address.ID.String())
				}
				return err
			}
			if err := s.insert(blocks); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Sync) insert(blocks []types.IBlock) error {
	for _, block := range blocks {
		select {
		case _, _ = <-s.stop:
			return nil
		default:
			if err := s.chain.Insert(block); err != nil {
				log.Warn("Insert chain failed!", "error", err, "height", block.GetHeight())
				if s.headerValidation(block.BlockHeader()) {
					s.fallBack()
					return err
				}
				s.peers.RemovePeer(s.curPeer.Address.ID.String())
				return err
			}
		}
	}
	return nil
}

// Remotely verify the block, if the block height is less than
// the effective block height, then discard the block. If the
// block occupies the majority of the currently started super
// nodes, it means that the block is more likely to be correct,
// and the block verification is successful.
func (s *Sync) headerValidation(header types.IHeader) bool {
	localEqual := false
	if header.GetHeight() <= s.chain.LastConfirmed() {
		return false
	}
	localHeader, err := s.chain.GetHeaderHeight(header.GetHeight())
	if err == nil && localHeader.GetHash().IsEqual(header.GetHash()) {
		localEqual = true
	}
	return s.validation(header, localEqual)

}

func (s *Sync) validation(header types.IHeader, localEqual bool) bool {
	count := 0
	ids := s.dPos.SuperIds()
	for _, id := range ids {
		if id != s.peers.Local().Address.ID.String() {
			peer := s.peers.Peer(id)
			if peer != nil {
				rs, err := s.request.IsEqual(peer.Conn, header)
				if err == nil && rs {
					count++
				}
			}
		} else if localEqual {
			count++
		}
	}
	if count > len(ids)/2 {
		return true
	}
	return false
}

// Block chain rolls back to a valid block
func (s *Sync) fallBack() {
	s.chain.Roll()
}

// Process blocks received from other super nodes.If the height
// of the block is greater than the local height, the storage is
// directly verified. If the height is less than the local height,
// the remote verification is performed, and the verification is
// passed back to the local block.
func (s *Sync) ReceivedBlockFromPeer(block types.IBlock) error {
	localHeight := s.chain.LastHeight()
	if localHeight == block.GetHeight()-1 {
		if err := s.chain.Insert(block); err != nil {
			log.Warn("Failed to insert received block", "err", err, "height", block.GetHeight(), "singer", block.GetSigner().String())
			return err
		}
	} else if block.GetHeight() <= localHeight {
		if localHeader, err := s.chain.GetBlockHeight(block.GetHeight()); err == nil {
			if !localHeader.GetHash().IsEqual(block.GetHash()) {
				if s.headerValidation(block.BlockHeader()) {
					s.fallBack()
					return err
				} else {
					log.Warn("Remote validation failed!", "height", block.GetHeight(), "signer", block.GetSigner().String())
					return err
				}
			}
		} else {
			if err := s.chain.Insert(block); err != nil {
				log.Warn("Failed to insert received block", "err", err, "height", block.GetHeight(), "singer", block.GetSigner().String())
				return err
			}
		}
	}
	return nil
}
