package rpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Futuremine-chain/future/common/blockchain"
	"github.com/Futuremine-chain/future/common/config"
	"github.com/Futuremine-chain/future/common/param"
	"github.com/Futuremine-chain/future/common/status"
	"github.com/Futuremine-chain/future/future/common/kit"
	"github.com/Futuremine-chain/future/future/common/kit/message"
	rpctypes "github.com/Futuremine-chain/future/future/rpc/types"
	fmctypes "github.com/Futuremine-chain/future/future/types"
	"github.com/Futuremine-chain/future/service/peers"
	"github.com/Futuremine-chain/future/service/pool"
	"github.com/Futuremine-chain/future/tools/arry"
	"github.com/Futuremine-chain/future/tools/crypto/certgen"
	log "github.com/Futuremine-chain/future/tools/log/log15"
	"github.com/Futuremine-chain/future/tools/utils"
	"github.com/Futuremine-chain/future/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"strconv"
)

const module = "rpc"

type Rpc struct {
	grpcServer *grpc.Server
	status     status.IStatus
	msgPool    *pool.Pool
	chain      blockchain.IChain
	peers      *peers.Peers
	getLocal   func() *types.Local
}

func NewRpc(status status.IStatus, msgPool *pool.Pool, chain blockchain.IChain, peers *peers.Peers) *Rpc {
	return &Rpc{status: status, msgPool: msgPool, chain: chain, peers: peers}
}

func (r *Rpc) Name() string {
	return module
}

func (r *Rpc) Start() error {
	lis, err := net.Listen("tcp", ":"+config.Param.RpcPort)
	if err != nil {
		return err
	}
	r.grpcServer, err = r.NewGRpcServer()
	if err != nil {
		return err
	}

	RegisterGreeterServer(r.grpcServer, r)
	reflection.Register(r.grpcServer)
	go func() {
		if err := r.grpcServer.Serve(lis); err != nil {
			log.Info("Rpc startup failed!", "module", module, "err", err)
			os.Exit(1)
			return
		}

	}()
	if config.Param.RpcTLS {
		log.Info("Rpc startup", "module", module, "port", config.Param.RpcPort, "pem", config.Param.RpcCert)
	} else {
		log.Info("Rpc startup", "module", module, "port", config.Param.RpcPort)
	}
	return nil
}

func (r *Rpc) Stop() error {
	r.grpcServer.Stop()
	log.Info("Rpc was stopped", "module", module)
	return nil
}

func (r *Rpc) Info() map[string]interface{} {
	return make(map[string]interface{}, 0)
}

func (r *Rpc) NewGRpcServer() (*grpc.Server, error) {
	var opts []grpc.ServerOption
	var interceptor grpc.UnaryServerInterceptor
	interceptor = r.interceptor
	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	// If tls is configured, generate tls certificate
	if config.Param.RpcTLS {
		if err := r.certFile(); err != nil {
			return nil, err
		}
		transportCredentials, err := credentials.NewServerTLSFromFile(config.Param.RpcCert, config.Param.RpcCertKey)
		if err != nil {
			return nil, err
		}
		opts = append(opts, grpc.Creds(transportCredentials))

	}

	// Set the maximum number of bytes received and sent
	opts = append(opts, grpc.MaxRecvMsgSize(param.MaxReqBytes))
	opts = append(opts, grpc.MaxSendMsgSize(param.MaxReqBytes))
	return grpc.NewServer(opts...), nil
}

func (r *Rpc) RegisterLocalInfo(f func() *types.Local) {
	r.getLocal = f
}

func (r *Rpc) GetAccount(_ context.Context, address *AddressReq) (*Response, error) {
	arryAddr := arry.StringToAddress(address.Address)
	if !kit.CheckAddress(config.Param.Name, arryAddr.String()) {
		return NewResponse(Err_Params, nil, fmt.Sprintf("%s address check failed", address.Address)), nil
	}
	account := r.status.Account(arryAddr)

	bytes, _ := json.Marshal(rpctypes.ToRpcAccount(account.(*fmctypes.Account)))
	return NewResponse(Success, bytes, ""), nil
}

func (r *Rpc) SendMessageRaw(ctx context.Context, code *SendMessageCodeReq) (*Response, error) {
	var rpcMsg *fmctypes.RpcMessage
	if err := json.Unmarshal(code.Code, &rpcMsg); err != nil {
		return NewResponse(Err_Params, nil, err.Error()), nil
	}
	tx, err := fmctypes.RpcMsgToMsg(rpcMsg)
	if err != nil {
		return NewResponse(Err_Params, nil, err.Error()), nil
	}
	if err := r.msgPool.Put(tx, false); err != nil {
		return NewResponse(Err_MsgPool, nil, err.Error()), nil
	}
	return NewResponse(Success, []byte(fmt.Sprintf("send message raw %s success", tx.Hash().String())), ""), nil
}

func (r *Rpc) GetMessage(ctx context.Context, hash *HashReq) (*Response, error) {
	hashArry, err := arry.StringToHash(hash.Hash)
	if err != nil {
		return NewResponse(Err_Params, nil, "wrong hash "+err.Error()), nil
	}
	msg, err := r.chain.GetMessage(hashArry)
	if err != nil {
		return NewResponse(Err_Chain, nil, err.Error()), nil
	}
	index, err := r.chain.GetMessageIndex(hashArry)
	if err != nil {
		return NewResponse(Err_Chain, nil, fmt.Sprintf("%s is not exist", hash.Hash)), nil
	}
	confirmed := r.chain.LastConfirmed()
	height := index.GetHeight()
	rpcMsg, _ := fmctypes.MsgToRpcMsg(msg.(*fmctypes.Message))
	rsMsg := &fmctypes.RpcMessageWithHeight{
		MsgHeader: rpcMsg.MsgHeader,
		MsgBody:   rpcMsg.MsgBody,
		Height:    height,
		Confirmed: confirmed >= height,
	}
	bytes, _ := json.Marshal(rsMsg)

	return NewResponse(Success, bytes, ""), nil
}

func (r *Rpc) GetBlockHash(ctx context.Context, hash *HashReq) (*Response, error) {
	hashArry, err := arry.StringToHash(hash.Hash)
	if err != nil {
		return NewResponse(Err_Params, nil, "wrong hash"), nil
	}
	block, err := r.chain.GetBlockHash(hashArry)
	if err != nil {
		return NewResponse(Err_Chain, nil, err.Error()), nil
	}
	rpcBlock, err := rpctypes.BlockToRpcBlock(block.(*fmctypes.Block), r.chain.LastConfirmed())
	if err != nil {
		return NewResponse(Err_Chain, nil, err.Error()), nil
	}
	bytes, _ := json.Marshal(rpcBlock)
	return NewResponse(Success, bytes, ""), nil
}

func (r *Rpc) GetBlockHeight(ctx context.Context, height *HeightReq) (*Response, error) {
	block, err := r.chain.GetBlockHeight(height.Height)
	if err != nil {
		return NewResponse(Err_Chain, nil, err.Error()), nil
	}
	rpcBlock, err := rpctypes.BlockToRpcBlock(block.(*fmctypes.Block), r.chain.LastConfirmed())
	if err != nil {
		return NewResponse(Err_Chain, nil, err.Error()), nil
	}
	bytes, _ := json.Marshal(rpcBlock)
	return NewResponse(Success, bytes, ""), nil
}

func (r *Rpc) LastHeight(context.Context, *NullReq) (*Response, error) {
	height := r.chain.LastHeight()
	sHeight := strconv.FormatUint(height, 10)
	return NewResponse(Success, []byte(sHeight), ""), nil
}
func (r *Rpc) Confirmed(context.Context, *NullReq) (*Response, error) {
	height := r.chain.LastConfirmed()
	sHeight := strconv.FormatUint(height, 10)
	return NewResponse(Success, []byte(sHeight), ""), nil
}

func (r *Rpc) GetMsgPool(context.Context, *NullReq) (*Response, error) {
	preparedTxs, futureTxs := r.msgPool.All()
	txPoolTxs := rpctypes.MsgsToRpcMsgsPool(preparedTxs, futureTxs)
	bytes, _ := json.Marshal(txPoolTxs)
	return NewResponse(Success, bytes, ""), nil
}

func (r *Rpc) Candidates(context.Context, *NullReq) (*Response, error) {
	candidates := r.status.Candidates()
	if candidates == nil || candidates.Len() == 0 {
		return NewResponse(Err_DPos, nil, "no candidates"), nil
	}
	cas := candidates.(*fmctypes.Candidates)
	for i, can := range cas.Members {
		for _, v := range can.Voters {
			cas.Members[i].Weight += r.chain.Vote(v)
		}
	}
	bytes, _ := json.Marshal(rpctypes.CandidatesToRpcCandidates(cas))
	return NewResponse(Success, bytes, ""), nil
}

func (r *Rpc) GetCycleSupers(ctx context.Context, cycle *CycleReq) (*Response, error) {
	supers := r.status.CycleSupers(cycle.Cycle)
	if supers == nil {
		return NewResponse(Err_DPos, nil, "no supers"), nil
	}
	bytes, _ := json.Marshal(rpctypes.SupersToRpcCandidates(supers))

	return NewResponse(Success, bytes, ""), nil
}

func (r *Rpc) Token(ctx context.Context, token *TokenAddressReq) (*Response, error) {
	iToken, err := r.status.Token(arry.StringToAddress(token.Token))
	if err != nil {
		return NewResponse(Err_Token, nil, fmt.Sprintf("token address %s is not exist", token.Token)), nil
	}
	bytes, _ := json.Marshal(rpctypes.TokenToRpcToken(iToken.(*fmctypes.TokenRecord)))
	return NewResponse(Success, bytes, ""), nil
}

func (r *Rpc) PeersInfo(context.Context, *NullReq) (*Response, error) {
	peersInfo := r.peers.PeersInfo()
	bytes, _ := json.Marshal(peersInfo)
	return NewResponse(Success, bytes, ""), nil
}
func (r *Rpc) LocalInfo(context.Context, *NullReq) (*Response, error) {
	if r.getLocal != nil {
		local := r.getLocal()
		bytes, _ := json.Marshal(local)
		return NewResponse(Success, bytes, ""), nil
	}
	return NewResponse(Err_Local, nil, "no local info"), nil
}

func (r *Rpc) GenerateAddress(ctx context.Context, req *GenerateReq) (*Response, error) {
	address, err := kit.GenerateAddress(req.Network, req.Publickey)
	if err != nil {
		return NewResponse(Err_Unknown, nil, err.Error()), nil
	}
	return NewResponse(Success, []byte(address), ""), nil
}

func (r *Rpc) GenerateTokenAddress(ctx context.Context, req *GenerateTokenReq) (*Response, error) {
	address, err := kit.GenerateTokenAddress(req.Network, req.Address, req.Abbr)
	if err != nil {
		return NewResponse(Err_Unknown, nil, err.Error()), nil
	}
	return NewResponse(Success, []byte(address), ""), nil
}

func (r *Rpc) CreateTransaction(ctx context.Context, req *TransactionReq) (*Response, error) {
	message := message.NewTransaction(req.From, req.To, req.Token, req.Amount, req.Fees, req.Nonce, req.Timestamp)
	bytes, _ := json.Marshal(message)
	return NewResponse(Success, bytes, ""), nil
}

func (r *Rpc) CreateToken(ctx context.Context, req *TokenReq) (*Response, error) {
	message := message.NewToken(req.From, req.Receiver, req.Token, req.Amount, req.Fees, req.Nonce, req.Timestamp, req.Name, req.Abbr, req.Increase)
	bytes, _ := json.Marshal(message)
	return NewResponse(Success, bytes, ""), nil
}

func (r *Rpc) CreateCandidate(ctx context.Context, req *CandidateReq) (*Response, error) {
	message := message.NewCandidate(req.From, req.P2Pid, req.Fees, req.Nonce, req.Timestamp)
	bytes, _ := json.Marshal(message)
	return NewResponse(Success, bytes, ""), nil
}
func (r *Rpc) CreateCancel(ctx context.Context, req *CancelReq) (*Response, error) {
	message := message.NewCancel(req.From, req.Fees, req.Nonce, req.Timestamp)
	bytes, _ := json.Marshal(message)
	return NewResponse(Success, bytes, ""), nil
}
func (r *Rpc) CreateVote(ctx context.Context, req *VoteReq) (*Response, error) {
	message := message.NewVote(req.From, req.To, req.Fees, req.Nonce, req.Timestamp)
	bytes, _ := json.Marshal(message)
	return NewResponse(Success, bytes, ""), nil
}

func (r *Rpc) SendTransaction(ctx context.Context, req *TransactionReq) (*Response, error) {
	message := message.NewTransaction(req.From, req.To, req.Token, req.Amount, req.Fees, req.Nonce, req.Timestamp)

	signature, err := hex.DecodeString(req.Signature)
	if err != nil {
		return NewResponse(Err_Params, nil, err.Error()), nil
	}
	pubKey, err := hex.DecodeString(req.Publickey)
	if err != nil {
		return NewResponse(Err_Params, nil, err.Error()), nil
	}
	message.Header.Signature.Bytes = signature
	message.Header.Signature.PubKey = pubKey

	if err := r.msgPool.Put(message, false); err != nil {
		return NewResponse(Err_MsgPool, nil, err.Error()), nil
	}
	return NewResponse(Success, []byte(fmt.Sprintf("send message %s success", message.Hash().String())), ""), nil
}

func (r *Rpc) SendToken(ctx context.Context, req *TokenReq) (*Response, error) {
	message := message.NewToken(req.From, req.Receiver, req.Token, req.Amount, req.Fees, req.Nonce, req.Timestamp, req.Name, req.Abbr, req.Increase)

	fmt.Println()
	x, _ := json.Marshal(message)
	fmt.Println(string(x))
	signature, err := hex.DecodeString(req.Signature)
	if err != nil {
		return NewResponse(Err_Params, nil, err.Error()), nil
	}
	pubKey, err := hex.DecodeString(req.Publickey)
	if err != nil {
		return NewResponse(Err_Params, nil, err.Error()), nil
	}
	message.Header.Signature.Bytes = signature
	message.Header.Signature.PubKey = pubKey

	if err := r.msgPool.Put(message, false); err != nil {
		return NewResponse(Err_MsgPool, nil, err.Error()), nil
	}
	return NewResponse(Success, []byte(fmt.Sprintf("send message %s success", message.Hash().String())), ""), nil
}

func (r *Rpc) SendCandidate(ctx context.Context, req *CandidateReq) (*Response, error) {
	message := message.NewCandidate(req.From, req.P2Pid, req.Fees, req.Nonce, req.Timestamp)
	signature, err := hex.DecodeString(req.Signature)
	if err != nil {
		return NewResponse(Err_Params, nil, err.Error()), nil
	}
	pubKey, err := hex.DecodeString(req.Publickey)
	if err != nil {
		return NewResponse(Err_Params, nil, err.Error()), nil
	}
	message.Header.Signature.Bytes = signature
	message.Header.Signature.PubKey = pubKey

	if err := r.msgPool.Put(message, false); err != nil {
		return NewResponse(Err_MsgPool, nil, err.Error()), nil
	}
	return NewResponse(Success, []byte(fmt.Sprintf("send message %s success", message.Hash().String())), ""), nil
}

func (r *Rpc) SendCancel(ctx context.Context, req *CancelReq) (*Response, error) {
	message := message.NewCancel(req.From, req.Fees, req.Nonce, req.Timestamp)
	signature, err := hex.DecodeString(req.Signature)
	if err != nil {
		return NewResponse(Err_Params, nil, err.Error()), nil
	}
	pubKey, err := hex.DecodeString(req.Publickey)
	if err != nil {
		return NewResponse(Err_Params, nil, err.Error()), nil
	}
	message.Header.Signature.Bytes = signature
	message.Header.Signature.PubKey = pubKey

	if err := r.msgPool.Put(message, false); err != nil {
		return NewResponse(Err_MsgPool, nil, err.Error()), nil
	}
	return NewResponse(Success, []byte(fmt.Sprintf("send message %s success", message.Hash().String())), ""), nil
}

func (r *Rpc) SendVote(ctx context.Context, req *VoteReq) (*Response, error) {
	message := message.NewVote(req.From, req.To, req.Fees, req.Nonce, req.Timestamp)
	signature, err := hex.DecodeString(req.Signature)
	if err != nil {
		return NewResponse(Err_Params, nil, err.Error()), nil
	}
	pubKey, err := hex.DecodeString(req.Publickey)
	if err != nil {
		return NewResponse(Err_Params, nil, err.Error()), nil
	}
	message.Header.Signature.Bytes = signature
	message.Header.Signature.PubKey = pubKey

	if err := r.msgPool.Put(message, false); err != nil {
		return NewResponse(Err_MsgPool, nil, err.Error()), nil
	}
	return NewResponse(Success, []byte(fmt.Sprintf("send message %s success", message.Hash().String())), ""), nil
}

func NewResponse(code int32, result []byte, err string) *Response {
	return &Response{Code: code, Result: result, Err: err}
}

// Authenticate rpc users
func (r *Rpc) auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("no token authentication information")
	}
	var (
		password string
	)

	if val, ok := md["password"]; ok {
		password = val[0]
	}

	if password != config.Param.RpcPass {
		return fmt.Errorf("the token authentication information is invalid: password=%s", password)
	}
	return nil
}

func (r *Rpc) interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	err = r.auth(ctx)
	if err != nil {
		return
	}
	return handler(ctx, req)
}

func (r *Rpc) certFile() error {
	if config.Param.RpcCert == "" {
		config.Param.RpcCert = config.Param.Data + "/server.pem"
	}
	if config.Param.RpcCertKey == "" {
		config.Param.RpcCertKey = config.Param.Data + "/server.key"
	}
	if !utils.Exist(config.Param.RpcCert) || !utils.Exist(config.Param.RpcCertKey) {
		return certgen.GenCertPair(config.Param.RpcCert, config.Param.RpcCertKey)
	}
	return nil
}
