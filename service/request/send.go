package request

import (
	"fmt"
	"github.com/Futuremine-chain/futuremine/service/peers"
	"github.com/Futuremine-chain/futuremine/tools/rlp"
	"github.com/Futuremine-chain/futuremine/tools/utils"
	"time"
)

var (
	lastHeight = Method("lastHeight")
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
