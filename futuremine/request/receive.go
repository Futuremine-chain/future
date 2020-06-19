package request

import (
	"github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/rlp"
)

/*lastHeight = Method("lastHeight")
sendTx     = Method("sendTx")
sendBlock  = Method("sendBlock")
getBlocks  = Method("getBlocks")
getBlock   = Method("getBlock")
isEqual    = Method("isEqual")*/

const maxSyncCount = 1000

func (r *RequestHandler) respLastHeight(req *ReqStream) (*Response, error) {
	var message string
	var body []byte
	code := Success
	height := r.chain.LastHeight()
	body, err := rlp.EncodeToBytes(height)
	if err != nil {
		code = Failed
		message = err.Error()
	}
	response := NewResponse(code, message, body)
	return response, nil
}

func (r *RequestHandler) respSendTx(req *ReqStream) (*Response, error) {
	defer func(){
		req.stream.Reset()
		req.stream.Close()
	}()

	var message string
	var body []byte
	code := Success
	msg, err := types.DecodeMessage(req.request.Body)
	if err != nil {
		code = Failed
		message = err.Error()
	}
	r.receiveMessage(msg.ToMessage())
	response := NewResponse(code, message, body)
	return response, nil
}

func (r *RequestHandler) respSendBlock(req *ReqStream) (*Response, error) {
	var message string
	var body []byte
	code := Success
	rlpBlock, err := types.DecodeRlpBlock(req.request.Body)
	if err != nil {
		code = Failed
		message = err.Error()
	}
	r.receiveBlock(rlpBlock.ToBlock())
	response := NewResponse(code, message, body)
	return response, nil
}

func (r *RequestHandler) respGetBlocks(req *ReqStream) (*Response, error) {
	var message string
	var body []byte
	var height uint64
	var count uint64
	code := Success
	rlpBlocks := make(types.RlpBlocks, 0)
	lastHeight := r.chain.LastHeight()
	err := rlp.DecodeBytes(req.request.Body, &height)
	if err != nil {
		code = Failed
		message = err.Error()
	} else if lastHeight >= height {
		for lastHeight >= height && count < maxSyncCount {
			block, err := r.chain.GetRlpBlockHeight(height)
			if err != nil {
				code = Failed
				message = err.Error()
				response := NewResponse(code, message, body)
				return response, nil
			} else {
				rlpBlocks = append(rlpBlocks, block)
				height++
				count++
			}
		}
		body, _ = rlpBlocks.Encode()
	} else {
		code = Failed
		message = ""
	}

	response := NewResponse(code, message, body)
	return response, nil
}
