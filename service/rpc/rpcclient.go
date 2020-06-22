package rpc

import (
	"context"
	"fmt"
	"github.com/Futuremine-chain/futuremine/common/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const timeout = 30

type Client struct {
	conn *grpc.ClientConn
	Gc   GreeterClient
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect() error {
	var conn *grpc.ClientConn
	var err error
	var opts []grpc.DialOption
	if config.Param.RpcTLS {
		creds, err := credentials.NewClientTLSFromFile(config.Param.RpcCert, "")
		if err != nil {
			return fmt.Errorf("failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(opts, grpc.WithPerRPCCredentials(&customCredential{Password: config.Param.RpcPass, OpenTLS: config.Param.RpcTLS}))

	conn, err = grpc.Dial(config.Param.RpcIp+":"+config.Param.RpcPort, opts...)
	if err != nil {
		return err
	}

	c.conn = conn
	c.Gc = NewGreeterClient(c.conn)
	return nil
}

func (c *Client) Close() {
	_ = c.conn.Close()
}

type customCredential struct {
	OpenTLS  bool
	Password string
}

func (c *customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"password": c.Password,
	}, nil
}

func (c *customCredential) RequireTransportSecurity() bool {
	if c.OpenTLS {
		return true
	}

	return false
}
