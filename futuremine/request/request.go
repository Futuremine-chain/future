package request

import (
	"encoding/json"
	"fmt"
	"github.com/Futuremine-chain/futuremine/common/blockchain"
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/common/param"
	log "github.com/Futuremine-chain/futuremine/tools/log/log15"
	"github.com/Futuremine-chain/futuremine/types"
	"github.com/libp2p/go-libp2p-core/network"
	"sync"
	"time"
)

const module = "request"

type request func(*ReqStream) (*Response, error)

type ReqStream struct {
	request *Request
	stream  network.Stream
}

func NewReqStream(r *Request, stream network.Stream) *ReqStream {
	return &ReqStream{r, stream}
}

func (r *ReqStream) Close() {
	r.stream.Reset()
	r.stream.Close()
}

type RequestHandler struct {
	chain          blockchain.IChain
	readyCh        chan *ReqStream
	bytesPool      sync.Pool
	receiveBlock   func(block types.IBlock) error
	receiveMessage func(msg types.IMessage) error
}

func NewRequestHandler(chain blockchain.IChain) *RequestHandler {
	return &RequestHandler{
		chain:   chain,
		readyCh: make(chan *ReqStream, config.Param.PeerRequestChan),
		bytesPool: sync.Pool{
			New: func() interface{} {
				return make([]byte, param.MaxReadBytes)
			},
		},
	}
}

func (r *RequestHandler) Name() string {
	return module
}

func (r *RequestHandler) Stop() error {
	return nil
}

// Listen for message requests
func (r *RequestHandler) Start() error {
	log.Info("Request handler started successfully", "module", module)
	go r.dealRequest()
	return nil
}

func (r *RequestHandler) dealRequest() {
	var h handler
	for reqStream := range r.readyCh {
		switch reqStream.request.Method {
		case sendBlock:
			h = r.respSendBlock
		default:
			reqStream.stream.Reset()
			reqStream.stream.Close()
			continue
		}
		go response(reqStream, h)
	}
}

func (r *RequestHandler) RegisterReceiveBlock(f func(types.IBlock) error) {
	r.receiveBlock = f
}

func (r *RequestHandler) RegisterReceiveMessage(f func(types.IMessage) error) {
	r.receiveMessage = f
}

// Handling message requests
func response(req *ReqStream, h handler) {
	defer req.Close()

	if response, err := h(req); err != nil {
		log.Warn("Response error", "module", module,
			"method", req.request.Method,
			"peer", req.stream.Conn().RemotePeer(),
			"addr", req.stream.Conn().RemoteMultiaddr(),
			"error", err)
	} else if response != nil {
		if err := responseStream(response, req.stream); err != nil {
			log.Warn("Send response error", "module", module,
				"method", req.request.Method,
				"peer", req.stream.Conn().RemotePeer(),
				"addr", req.stream.Conn().RemoteMultiaddr(),
				"error", err)
		}
	}
	for isAlive(req.stream) {
		time.Sleep(time.Second * 1)
	}
}

func (r *RequestHandler) SendToReady(stream network.Stream) {
	request, err := r.UnmarshalRequest(stream)
	if err != nil {
		return
	}
	r.readyCh <- NewReqStream(request, stream)
}

// Read from request
func (r *RequestHandler) UnmarshalRequest(stream network.Stream) (*Request, error) {
	reBytes, err := r.read(stream)
	request := &Request{}

	err = json.Unmarshal(reBytes, request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

// Read from response
func (r *RequestHandler) UnmarshalResponse(stream network.Stream) (*Response, error) {
	reBytes, err := r.read(stream)
	if err != nil {
		fmt.Println("UnmarshalResponse", err)
		return nil, err
	}
	resp := &Response{}

	err = json.Unmarshal(reBytes, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Read message bytes
func (rm *RequestHandler) read(stream network.Stream) ([]byte, error) {
	var rs []byte
	var err error
	var n int
	len := 0
	arry := rm.bytesPool.Get().([]byte)
	for len < param.MaxReqBytes {
		reset(arry)
		n, err = stream.Read(arry)
		if err != nil {
			break
		}
		rs = append(rs, arry...)
		len += n
		if n < param.MaxReadBytes {
			break
		}
	}
	rm.bytesPool.Put(arry)
	if len > param.MaxReqBytes {
		return nil, fmt.Errorf("request data must be less than %d", param.MaxReqBytes)
	}
	return rs[0:len], err
}

type handler func(*ReqStream) (*Response, error)

func reset(bytes []byte) {
	for i, _ := range bytes {
		bytes[i] = 0
	}
}

func responseStream(response *Response, stream network.Stream) error {
	bytes, err := json.Marshal(response)
	if err != nil {
		return err
	}
	_, err = stream.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func requestStream(request *Request, stream network.Stream) error {
	bytes, err := json.Marshal(request)
	if err != nil {
		return err
	}
	_, err = stream.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func isAlive(stream network.Stream) bool {
	bytes := [10]byte{}
	_, err := stream.Read(bytes[:])
	if err != nil {
		return false
	}
	return true
}
