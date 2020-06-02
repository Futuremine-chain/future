package request

import (
	"github.com/Futuremine-chain/futuremine/tools/rlp"
)

func (r *RequestHandler) respLastHeight(req *ReqStream) (*Response, error) {
	var message string
	var body []byte
	code := Success
	height := r.chain.LastHeight()
	body, _ = rlp.EncodeToBytes(height)
	response := NewResponse(code, message, body)
	return response, nil
}
