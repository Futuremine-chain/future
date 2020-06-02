package request

import (
	"fmt"
	fmctypes "github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/service/peers"
	request2 "github.com/Futuremine-chain/futuremine/service/request"
	"github.com/Futuremine-chain/futuremine/tools/rlp"
	"github.com/Futuremine-chain/futuremine/tools/utils"
	"github.com/Futuremine-chain/futuremine/types"
	"time"
)

var (
	lastHeight = Method("lastHeight")
	sendTx     = Method("sendTx")
	sendBlock  = Method("sendBlock")
	getBlocks  = Method("getBlocks")
	getBlock   = Method("getBlock")
)

func (r *RequestHandler) LastHeight(conn *peers.Conn) (uint64, error) {
	var height uint64 = 0
	s, err := conn.Stream.Conn().NewStream()
	if err != nil {
		return 0, err
	}
	defer func() {
		s.Reset()
		s.Close()
	}()

	s.SetDeadline(time.Unix(utils.NowUnix()+timeOut, 0))
	req := NewRequest(lastHeight, nil)
	err = requestStream(req, s)
	if err != nil {
		return 0, err
	}
	response, err := r.UnmarshalResponse(s)
	if response != nil && response.Code == Success {
		err := rlp.DecodeBytes(response.Body, &height)
		if err != nil {
			return 0, err
		}
	} else {
		return 0, fmt.Errorf("peer error: %v", err)
	}
	return height, nil
}

func (r *RequestHandler) SendTx(conn *peers.Conn, tx types.ITransaction) error {
	s, err := conn.Stream.Conn().NewStream()
	if err != nil {
		return err
	}
	defer func() {
		s.Reset()
		s.Close()
	}()

	s.SetDeadline(time.Unix(utils.NowUnix()+timeOut, 0))
	req := NewRequest(sendTx, tx.ToRlp().Bytes())
	err = requestStream(req, s)
	if err != nil {
		return err
	}
	response, err := r.UnmarshalResponse(s)
	if response != nil && response.Code == Success {
		return nil
	} else {
		return fmt.Errorf("peer error: %v", err)
	}
}

func (r *RequestHandler) SendBlock(conn *peers.Conn, block types.IBlock) error {
	s, err := conn.Stream.Conn().NewStream()
	if err != nil {
		return err
	}
	defer func() {
		s.Reset()
		s.Close()
	}()

	s.SetDeadline(time.Unix(utils.NowUnix()+timeOut, 0))
	//body := xx
	req := NewRequest(sendBlock, block.ToRlp().Bytes())
	err = requestStream(req, s)
	if err != nil {
		return err
	}
	response, err := r.UnmarshalResponse(s)
	if response != nil && response.Code == Success {
		return nil
	} else {
		return fmt.Errorf("peer error: %v", err)
	}
}

func (r *RequestHandler) GetBlocks(conn *peers.Conn, height uint64) ([]types.IBlock, error) {
	s, err := conn.Stream.Conn().NewStream()
	if err != nil {
		return nil, err
	}
	defer func() {
		s.Reset()
		s.Close()
	}()

	bytes, err := rlp.EncodeToBytes(height)
	if err != nil {
		return nil, err
	}
	s.SetDeadline(time.Unix(utils.NowUnix()+timeOut, 0))
	request := NewRequest(getBlocks, bytes)
	err = requestStream(request, s)
	if err != nil {
		return nil, request2.Err_PeerClosed
	}
	response, err := r.UnmarshalResponse(s)
	if response != nil && response.Code == Success {
		blocks, err := fmctypes.DecodeRlpBlocks(response.Body)
		if err != nil {
			return nil, err
		}
		return blocks.ToBlocks(), nil
	} else if response != nil && response.Message == request2.Err_BlockNotFound.Error() {
		return nil, request2.Err_BlockNotFound
	} else {
		return nil, request2.Err_PeerClosed
	}
}

func (r *RequestHandler) GetBlock(conn *peers.Conn, height uint64) (types.IBlock, error) {
	s, err := conn.Stream.Conn().NewStream()
	if err != nil {
		return nil, err
	}
	defer func() {
		s.Reset()
		s.Close()
	}()

	bytes, err := rlp.EncodeToBytes(height)
	if err != nil {
		return nil, err
	}
	s.SetDeadline(time.Unix(utils.NowUnix()+timeOut, 0))
	request := NewRequest(getBlocks, bytes)
	err = requestStream(request, s)
	if err != nil {
		return nil, request2.Err_PeerClosed
	}
	response, err := r.UnmarshalResponse(s)
	if response != nil && response.Code == Success {
		block, err := fmctypes.DecodeRlpBlock(response.Body)
		if err != nil {
			return nil, err
		}
		return block.ToBlock(), nil
	} else if response != nil && response.Message == request2.Err_BlockNotFound.Error() {
		return nil, request2.Err_BlockNotFound
	} else {
		return nil, request2.Err_PeerClosed
	}
}
