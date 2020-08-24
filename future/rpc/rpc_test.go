package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Futuremine-chain/future/common/config"
	"github.com/Futuremine-chain/future/future/types"
	"github.com/Futuremine-chain/future/tools/arry"
	"github.com/Futuremine-chain/future/tools/crypto/ecc/secp256k1"
	"testing"
	"time"
)

var client = NewClient(&config.RpcConfig{
	DataDir:    "",
	RpcIp:      "127.0.0.1",
	RpcPort:    "29163",
	RpcTLS:     false,
	RpcCert:    "",
	RpcCertKey: "",
	RpcPass:    "",
})

type Header struct {
	Hash string
}
type respHash struct {
	Header *Header
}

func TestRpc_GenerateAddress(t *testing.T) {
	err := client.Connect()
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	resp, err := client.Gc.GenerateAddress(context.Background(), &GenerateReq{
		Network:   "mainnet",
		Publickey: "021c39e5bae2894676b8c70d8ba25f84ef70ac59440fcd585640bda8d02646d4b5",
	})
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	if bytes.Compare(resp.Result, []byte("FMejc9bjiTeQzKQG9fSDPGdsRzzEdEQe6se")) != 0 {
		t.Fatal("error")
	}
}

func TestRpc_GenerateTokenAddress(t *testing.T) {
	err := client.Connect()
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	resp, err := client.Gc.GenerateTokenAddress(context.Background(), &GenerateTokenReq{
		Network: "mainnet",
		Address: "FMejc9bjiTeQzKQG9fSDPGdsRzzEdEQe6se",
		Abbr:    "ANBJ",
	})
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	if bytes.Compare(resp.Result, []byte("FTgeasx9fmkEiVu69xr56hC9c1QTv4rKM8e")) != 0 {
		t.Fatal("error")
	}
}

func TestRpc_CreateTransaction(t *testing.T) {
	err := client.Connect()
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	ti := uint64(time.Now().Unix())
	resp, err := client.Gc.CreateTransaction(context.Background(), &TransactionReq{
		From:      "FMejc9bjiTeQzKQG9fSDPGdsRzzEdEQe6se",
		To:        "FMegukTco2m1S9Y4ebXM9kVpQ6jqGGZBwWv",
		Token:     "FM",
		Amount:    1000000000,
		Fees:      1000000,
		Nonce:     3,
		Timestamp: ti,
	})
	var h *respHash

	err = json.Unmarshal(resp.Result, &h)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	priv, err := secp256k1.ParseStringToPrivate("68d01d8fe1d512f9038040f0e1d3b26a599513a2e6595322aae07060afae698c")
	hash, err := arry.StringToHash(h.Header.Hash)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	si, err := types.Sign(priv, hash)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	resp, err = client.Gc.SendTransaction(context.Background(), &TransactionReq{
		From:      "FMejc9bjiTeQzKQG9fSDPGdsRzzEdEQe6se",
		To:        "FMegukTco2m1S9Y4ebXM9kVpQ6jqGGZBwWv",
		Token:     "FM",
		Amount:    1000000000,
		Fees:      1000000,
		Nonce:     3,
		Timestamp: ti,
		Signature: si.SignatureString(),
		Publickey: si.PubKeyString(),
	})
	fmt.Println(resp, err)
}

func TestRpc_CreateToken(t *testing.T) {
	err := client.Connect()
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	ti := uint64(time.Now().Unix())
	resp, err := client.Gc.GenerateTokenAddress(context.Background(), &GenerateTokenReq{
		Network: "mainnet",
		Address: "FMejc9bjiTeQzKQG9fSDPGdsRzzEdEQe6se",
		Abbr:    "ANBJ",
	})

	token := string(resp.Result)
	resp, err = client.Gc.CreateToken(context.Background(), &TokenReq{
		From:      "FMejc9bjiTeQzKQG9fSDPGdsRzzEdEQe6se",
		Receiver:  "FMegukTco2m1S9Y4ebXM9kVpQ6jqGGZBwWv",
		Token:     token,
		Amount:    1000000000000,
		Fees:      1000000,
		Nonce:     4,
		Name:      "12121",
		Abbr:      "ANBJ",
		Increase:  true,
		Timestamp: ti,
	})
	fmt.Println(string(resp.Result))
	var h *respHash

	err = json.Unmarshal(resp.Result, &h)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	priv, err := secp256k1.ParseStringToPrivate("68d01d8fe1d512f9038040f0e1d3b26a599513a2e6595322aae07060afae698c")
	hash, err := arry.StringToHash(h.Header.Hash)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	fmt.Println(hash.String())
	si, err := types.Sign(priv, hash)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	resp, err = client.Gc.SendToken(context.Background(), &TokenReq{
		From:      "FMejc9bjiTeQzKQG9fSDPGdsRzzEdEQe6se",
		Receiver:  "FMegukTco2m1S9Y4ebXM9kVpQ6jqGGZBwWv",
		Token:     token,
		Amount:    1000000000000,
		Fees:      1000000,
		Nonce:     4,
		Name:      "12121",
		Abbr:      "ANBJ",
		Increase:  true,
		Timestamp: ti,
		Signature: si.SignatureString(),
		Publickey: si.PubKeyString(),
	})
	fmt.Println(resp, err)
}

func TestRpc_CreateCandidate(t *testing.T) {
	err := client.Connect()
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	ti := uint64(time.Now().Unix())
	resp, err := client.Gc.CreateCandidate(context.Background(), &CandidateReq{
		From:      "FMejc9bjiTeQzKQG9fSDPGdsRzzEdEQe6se",
		P2Pid:     "16Uiu2HAkwKrbmaz3WRPjdJZbEBDCj412auZPoBCr3cpDViztzcX6",
		Fees:      1000000,
		Nonce:     5,
		Timestamp: ti,
	})
	var h *respHash

	err = json.Unmarshal(resp.Result, &h)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	priv, err := secp256k1.ParseStringToPrivate("68d01d8fe1d512f9038040f0e1d3b26a599513a2e6595322aae07060afae698c")
	hash, err := arry.StringToHash(h.Header.Hash)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	si, err := types.Sign(priv, hash)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	resp, err = client.Gc.SendCandidate(context.Background(), &CandidateReq{
		From:      "FMejc9bjiTeQzKQG9fSDPGdsRzzEdEQe6se",
		P2Pid:     "16Uiu2HAkwKrbmaz3WRPjdJZbEBDCj412auZPoBCr3cpDViztzcX6",
		Fees:      1000000,
		Nonce:     5,
		Timestamp: ti,
		Signature: si.SignatureString(),
		Publickey: si.PubKeyString(),
	})
	fmt.Println(resp, err)
}

func TestRpc_CreateCancel(t *testing.T) {
	err := client.Connect()
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	ti := uint64(time.Now().Unix())
	resp, err := client.Gc.CreateCancel(context.Background(), &CancelReq{
		From:      "FMejc9bjiTeQzKQG9fSDPGdsRzzEdEQe6se",
		Fees:      1000000,
		Nonce:     6,
		Timestamp: ti,
	})
	var h *respHash

	err = json.Unmarshal(resp.Result, &h)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	priv, err := secp256k1.ParseStringToPrivate("68d01d8fe1d512f9038040f0e1d3b26a599513a2e6595322aae07060afae698c")
	hash, err := arry.StringToHash(h.Header.Hash)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	si, err := types.Sign(priv, hash)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	resp, err = client.Gc.SendCancel(context.Background(), &CancelReq{
		From:      "FMejc9bjiTeQzKQG9fSDPGdsRzzEdEQe6se",
		Fees:      1000000,
		Nonce:     6,
		Timestamp: ti,
		Signature: si.SignatureString(),
		Publickey: si.PubKeyString(),
	})
	fmt.Println(resp, err)
}

func TestRpc_CreateVote(t *testing.T) {
	err := client.Connect()
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	ti := uint64(time.Now().Unix())
	resp, err := client.Gc.CreateVote(context.Background(), &VoteReq{
		From:      "FMejc9bjiTeQzKQG9fSDPGdsRzzEdEQe6se",
		To:        "FMejc9bjiTeQzKQG9fSDPGdsRzzEdEQe6se",
		Fees:      1000000,
		Nonce:     7,
		Timestamp: ti,
	})
	var h *respHash

	err = json.Unmarshal(resp.Result, &h)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	priv, err := secp256k1.ParseStringToPrivate("68d01d8fe1d512f9038040f0e1d3b26a599513a2e6595322aae07060afae698c")
	hash, err := arry.StringToHash(h.Header.Hash)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	si, err := types.Sign(priv, hash)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	resp, err = client.Gc.SendVote(context.Background(), &VoteReq{
		From:      "FMejc9bjiTeQzKQG9fSDPGdsRzzEdEQe6se",
		To:        "FMejc9bjiTeQzKQG9fSDPGdsRzzEdEQe6se",
		Fees:      1000000,
		Nonce:     7,
		Timestamp: ti,
		Signature: si.SignatureString(),
		Publickey: si.PubKeyString(),
	})
	fmt.Println(resp, err)
}
