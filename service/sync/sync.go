package sync

import (
	"errors"
	"github.com/Futuremine-chain/futuremine/common/blockchain"
	"github.com/Futuremine-chain/futuremine/service/peers"
	"github.com/Futuremine-chain/futuremine/service/request"
	log "github.com/Futuremine-chain/futuremine/tools/log/log15"
	"github.com/Futuremine-chain/futuremine/types"
	"math/rand"
	"time"
)

const module = "sync"

type Sync struct {
	chain   blockchain.IBlockChain
	request request.IRequestHandler
	peers   *peers.Peers
	curPeer *peers.Peer
	stop    chan bool
	stopped chan bool
}

func NewSync(peers *peers.Peers, request request.IRequestHandler, chain blockchain.IBlockChain) *Sync {
	s := &Sync{
		chain:   chain,
		peers:   peers,
		request: request,
		stop:    make(chan bool),
		stopped: make(chan bool),
	}
	request.RegisterReceiveBlock(s.receivedBlock)
	return s
}

func (s *Sync) Name() string {
	return module
}

func (s *Sync) Start() error {
	go s.syncBlocks()
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
				log.Warn("Insert chain failed!", "error", err, "height", block.Height())
				if s.superValidation(block.Header()) {
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
func (s *Sync) superValidation(header types.IHeader) bool {
	if header.Height() <= s.chain.LastConfirmed() {
		return false
	}
	/*ids, err := s.consensus.GetWinnersPeerID(header.Time())
	if err != nil {
		return false
	}
	localHeader, err := s.chain.GetHeaderHeight(header.Height())
	if err == nil {
		if localHeader.Hash().IsEqual(header.Hash()) {
			return false
		}
	}
	compareMap := make(map[string][]string)
	for _, id := range ids {
		peerId := new(peer.ID)
		if id != bm.peerManager.LocalPeerInfo().AddrInfo.ID.String() {
			if err = peerId.UnmarshalText([]byte(id)); err == nil {
				streamCreator := p2p.StreamCreator{PeerId: *peerId, NewStreamFunc: bm.newStream.CreateStream}
				remoteHeader, err := bm.network.GetHeaderByHeight(&streamCreator, header.Height)
				if err != nil {
					continue
				}
				if _, ok := compareMap[remoteHeader.HashString()]; ok {
					compareMap[remoteHeader.HashString()] = append(compareMap[remoteHeader.HashString()], id)
				} else {
					compareMap[remoteHeader.HashString()] = []string{id}
				}
			}
		} else {
			localHeader, err := bm.blockChain.GetHeaderByHeight(header.Height)
			if err != nil {
				return true
			}
			if _, ok := compareMap[localHeader.HashString()]; ok {
				compareMap[localHeader.HashString()] = append(compareMap[localHeader.HashString()], id)
			} else {
				compareMap[localHeader.HashString()] = []string{id}
			}
		}
	}
	selectedHash := getEffectiveHash(compareMap)
	if header.HashString() != selectedHash {
		return false
	}*/
	return true
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
func (s *Sync) receivedBlock(block types.IBlock) error {
	localHeight := s.chain.LastHeight()
	if localHeight == block.Height()-1 {
		if err := s.chain.Insert(block); err != nil {
			log.Warn("Failed to insert received block", "err", err, "height", block.Height(), "singer", block.Signer().String())
			return err
		}
	} else if block.Height() <= localHeight {
		if localHeader, err := s.chain.GetBlockHeight(block.Height()); err == nil {
			if !localHeader.Hash().IsEqual(block.Hash()) {
				if s.superValidation(block.Header()) {
					s.fallBack()
					return err
				} else {
					log.Warn("Remote validation failed!", "height", block.Height(), "signer", block.Signer().String())
					return err
				}
			}
		} else {
			if err := s.chain.Insert(block); err != nil {
				log.Warn("Failed to insert received block", "err", err, "height", block.Height(), "singer", block.Signer().String())
				return err
			}
		}
	}
	return nil
}

func getEffectiveHash(compareMap map[string][]string) string {
	hashes := make([]string, 0)
	var maxCount int
	for h, peers := range compareMap {
		if len(peers) == maxCount {
			hashes = append(hashes, h)
		} else if len(peers) > maxCount {
			maxCount = len(peers)
			hashes = []string{h}
		}
	}
	if len(hashes) > 1 {
		rand.Intn(len(hashes))
		return hashes[rand.Intn(len(hashes))]
	}
	if len(hashes) == 0 {
		return ""
	}
	return hashes[0]
}
