package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Futuremine-chain/futuremine/common/blockchain"
	"github.com/Futuremine-chain/futuremine/common/config"
	"github.com/Futuremine-chain/futuremine/common/param"
	"github.com/Futuremine-chain/futuremine/common/status"
	"github.com/Futuremine-chain/futuremine/futuremine/common/kit"
	rpctypes "github.com/Futuremine-chain/futuremine/futuremine/rpc/types"
	fmctypes "github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/service/pool"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/crypto/certgen"
	log "github.com/Futuremine-chain/futuremine/tools/log/log15"
	"github.com/Futuremine-chain/futuremine/tools/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

const module = "rpc"

type Rpc struct {
	grpcServer *grpc.Server
	status     status.IStatus
	msgPool    *pool.Pool
	chain      blockchain.IChain
}

func NewRpc(status status.IStatus, msgPool *pool.Pool, chain blockchain.IChain) *Rpc {
	return &Rpc{status: status, msgPool: msgPool, chain: chain}
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

func (r *Rpc) GetAccount(_ context.Context, req *Request) (*Response, error) {
	params := make([]interface{}, 0)
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewResponse(Err_Params, nil, err.Error()), nil
	}
	if len(params) == 0 {
		return NewResponse(Err_Params, nil, "no address"), nil
	}
	if address, ok := params[0].(string); !ok {
		return NewResponse(Err_Params, nil, "address type error"), nil
	} else {
		arryAddr := arry.StringToAddress(address)
		if !kit.CheckAddress(config.Param.Name, arryAddr) {
			return NewResponse(Err_Params, nil, fmt.Sprintf("%s address check failed", string(req.Params))), nil
		}
		account := r.status.Account(arryAddr)

		bytes, _ := json.Marshal(rpctypes.ToRpcAccount(account.(*fmctypes.Account)))
		return NewResponse(Success, bytes, ""), nil
	}
}

func (r *Rpc) SendMessageRaw(ctx context.Context, req *Request) (*Response, error) {
	var rpcMsg *rpctypes.RpcMessage
	if err := json.Unmarshal(req.Params, &rpcMsg); err != nil {
		return NewResponse(Err_Params, nil, err.Error()), nil
	}
	tx, err := rpctypes.RpcMsgToMsg(rpcMsg)
	if err != nil {
		return NewResponse(Err_Params, nil, ""), nil
	}
	if err := r.msgPool.Put(tx, false); err != nil {
		return NewResponse(Err_MsgPool, nil, err.Error()), nil
	}
	return NewResponse(Success, []byte(fmt.Sprintf("send transaction raw %s success", tx.Hash().String())), ""), nil
}
func (r *Rpc) GetMessage(ctx context.Context, req *Request) (*Response, error) {
	params := make([]interface{}, 0)
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewResponse(Err_Params, nil, err.Error()), nil
	}
	if len(params) < 1 {
		return NewResponse(Err_Params, nil, "no hash"), nil
	}
	if hashStr, ok := params[0].(string); !ok {
		return NewResponse(Err_Params, nil, "only string hash is allowed"), nil
	} else {
		hash, err := arry.StringToHash(hashStr)
		if err != nil {
			return NewResponse(Err_Params, nil, "wrong hash "+err.Error()), nil
		}
		msg, err := r.chain.GetMessage(hash)
		if err != nil {
			return NewResponse(Err_Chain, nil, err.Error()), nil
		}
		rpcMsg, _ := rpctypes.MsgToRpcMsg(msg.(*fmctypes.Message))
		bytes, _ := json.Marshal(rpcMsg)

		return NewResponse(Success, bytes, ""), nil
	}
}
func (r *Rpc) GetBlockHash(ctx context.Context, req *Request) (*Response, error) {
	params := make([]interface{}, 0)
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewResponse(Err_Params, nil, err.Error()), nil
	}
	if len(params) < 1 {
		return NewResponse(Err_Params, nil, "no hash"), nil
	}
	if hashStr, ok := params[0].(string); !ok {
		return NewResponse(Err_Params, nil, "only string hash is allowed"), nil
	} else {
		hash, err := arry.StringToHash(hashStr)
		if err != nil {
			return NewResponse(Err_Params, nil, "wrong hash"), nil
		}
		block, err := r.chain.GetBlockHash(hash)
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
	return nil, nil
}

func (r *Rpc) GetBlockHeight(context.Context, *Request) (*Response, error) { return nil, nil }
func (r *Rpc) LastHeight(context.Context, *Request) (*Response, error)     { return nil, nil }
func (r *Rpc) Confirmed(context.Context, *Request) (*Response, error)      { return nil, nil }
func (r *Rpc) GetMsgPool(context.Context, *Request) (*Response, error)     { return nil, nil }
func (r *Rpc) Candidates(context.Context, *Request) (*Response, error)     { return nil, nil }
func (r *Rpc) GetCycleSupers(context.Context, *Request) (*Response, error) { return nil, nil }
func (r *Rpc) Token(context.Context, *Request) (*Response, error)          { return nil, nil }
func (r *Rpc) PeersInfo(context.Context, *Request) (*Response, error)      { return nil, nil }
func (r *Rpc) LocalInfo(context.Context, *Request) (*Response, error)      { return nil, nil }

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
