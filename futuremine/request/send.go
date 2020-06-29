package request

import (
	"fmt"
	fmctypes "github.com/Futuremine-chain/futuremine/futuremine/types"
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
	isEqual    = Method("isEqual")
	localInfo  = Method("localInfo")
)

func (r *RequestHandler) LastHeight(conn *types.Conn) (uint64, error) {
	var height uint64 = 0
	s, err := conn.Create(conn.PeerId)
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

func (r *RequestHandler) SendMsg(conn *types.Conn, msg types.IMessage) error {
	s, err := conn.Create(conn.PeerId)
	if err != nil {
		return err
	}

	defer func() {
		s.Reset()
		s.Close()
	}()

	s.SetDeadline(time.Unix(utils.NowUnix()+timeOut, 0))
	req := NewRequest(sendTx, msg.ToRlp().Bytes())
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

func (r *RequestHandler) SendBlock(conn *types.Conn, block types.IBlock) error {
	s, err := conn.Create(conn.PeerId)
	fmt.Println("send block", conn.PeerId)
	if err != nil {
		return err
	}

	defer func() {
		s.Reset()
		s.Close()
	}()

	s.SetDeadline(time.Unix(utils.NowUnix()+timeOut, 0))
	//body := xx
	req := NewRequest(sendBlock, block.ToRlpBlock().Bytes())
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

func (r *RequestHandler) GetBlocks(conn *types.Conn, height uint64) ([]types.IBlock, error) {
	s, err := conn.Create(conn.PeerId)
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
		return fmctypes.RlpBlocksToBlocks(blocks), nil
	} else if response != nil && response.Message == request2.Err_BlockNotFound.Error() {
		return nil, request2.Err_BlockNotFound
	} else {
		return nil, request2.Err_PeerClosed
	}
}

func (r *RequestHandler) GetBlock(conn *types.Conn, height uint64) (types.IBlock, error) {
	s, err := conn.Create(conn.PeerId)
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

func (r *RequestHandler) IsEqual(conn *types.Conn, header types.IHeader) (bool, error) {
	s, err := conn.Create(conn.PeerId)
	if err != nil {
		return false, err
	}

	defer func() {
		s.Reset()
		s.Close()
	}()

	s.SetDeadline(time.Unix(time.Now().Unix()+timeOut, 0))

	request := NewRequest(isEqual, header.Bytes())
	err = requestStream(request, s)
	response, err := r.UnmarshalResponse(s)
	var rs bool
	if response != nil && response.Code == Success {
		err := rlp.DecodeBytes(response.Body, &rs)
		if err != nil {
			return false, err
		}
	} else {
		return false, fmt.Errorf("peer error: %v", err)
	}
	return rs, nil
}

func (r *RequestHandler) LocalInfo(conn *types.Conn) (*types.Local, error) {
	s, err := conn.Create(conn.PeerId)
	if err != nil {
		return nil, err
	}

	defer func() {
		s.Reset()
		s.Close()
	}()

	s.SetDeadline(time.Unix(time.Now().Unix()+timeOut, 0))

	request := NewRequest(localInfo, nil)
	err = requestStream(request, s)
	response, err := r.UnmarshalResponse(s)
	var rs *types.Local
	if response != nil && response.Code == Success {
		err := rlp.DecodeBytes(response.Body, &rs)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("peer error: %v", err)
	}
	return rs, nil
}
