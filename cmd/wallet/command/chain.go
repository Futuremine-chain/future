package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Futuremine-chain/futuremine/futuremine/common/kit/message"
	"github.com/Futuremine-chain/futuremine/futuremine/rpc"
	rpctypes "github.com/Futuremine-chain/futuremine/futuremine/rpc/types"
	"github.com/Futuremine-chain/futuremine/futuremine/types"
	amount2 "github.com/Futuremine-chain/futuremine/tools/amount"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/crypto/ecc/secp256k1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

func init() {
	blockCmds := []*cobra.Command{
		LastHeightCmd,
		GetBlockCmd,
		GetMessageCmd,
		SendMessageCmd,
	}

	RootCmd.AddCommand(blockCmds...)
	RootSubCmdGroups["chain"] = blockCmds
}

var LastHeightCmd = &cobra.Command{
	Use:     "LastHeight",
	Short:   "LastHeight; Get last height of node;",
	Aliases: []string{"lastheight", "LH", "lh"},
	Example: `
	LastHeight 
	`,
	Args: cobra.MinimumNArgs(0),
	Run:  LastHeight,
}

func LastHeight(cmd *cobra.Command, args []string) {
	client, err := NewRpcClient()
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()
	resp, err := client.Gc.LastHeight(ctx, &rpc.Null{})
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	if resp.Code == 0 {
		output(string(resp.Result))
		return
	}
	outputRespError(cmd.Use, resp)
}

var GetBlockCmd = &cobra.Command{
	Use:     "GetBlock {height/hash};",
	Short:   "GetBlock {height/hash}; Get block by height or hash;",
	Aliases: []string{"getblock", "gb", "GB"},
	Example: `
	GetBlock 1 
	GetBlock 0x4e32b712330c0d4ee45f06017390c5d1d3c26d0e6c7be4ea9a5036bdb6c72a07 
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  GetBlock,
}

func GetBlock(cmd *cobra.Command, args []string) {
	height, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		GetBlockByHash(cmd, args)
		return
	}
	client, err := NewRpcClient()
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()

	resp, err := client.Gc.GetBlockHeight(ctx, &rpc.Height{Height: height})
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	if resp.Code == 0 {
		output(string(resp.Result))
		return
	}
	outputRespError(cmd.Use, resp)

}

func GetBlockByHash(cmd *cobra.Command, args []string) {
	var err error
	client, err := NewRpcClient()
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()


	resp, err := client.Gc.GetBlockHash(ctx, &rpc.Hash{Hash: args[0]})
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	if resp.Code == 0 {
		output(string(resp.Result))
		return
	}
	outputRespError(cmd.Use, resp)

}

var SendMessageCmd = &cobra.Command{
	Use:     "SendTransaction {from} {to} {token} {amount} {fees} {password} {nonce}; Send a transaction;",
	Aliases: []string{"sendtransaction", "ST", "st"},
	Short:   "SendTransaction {from} {to} {token} {amount} {fees} {password} {nonce}; Send a transaction;",
	Example: `
	SendTransaction xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ xCE9boXz2TxSE9srVPDdfszyiXtfT3vduc8 FMC 10 0.1
		OR
	SendTransaction xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ xCE9boXz2TxSE9srVPDdfszyiXtfT3vduc8 FMC 10 0.1 123456
		OR
	SendTransaction xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ xCE9boXz2TxSE9srVPDdfszyiXtfT3vduc8 FMC 10 0.1 123456 1
	`,
	Args: cobra.MinimumNArgs(5),
	Run:  SendTransaction,
}

func SendTransaction(cmd *cobra.Command, args []string) {
	var passwd []byte
	var err error
	if len(args) > 5 {
		passwd = []byte(args[5])
	} else {
		fmt.Println("please input password：")
		passwd, err = readPassWd()
		if err != nil {
			log.Error(cmd.Use+" err: ", fmt.Errorf("read password failed! %s", err.Error()))
			return
		}
	}
	privKey, err := loadPrivate(getAddJsonPath(args[0]), passwd)
	if err != nil {
		log.Error(cmd.Use+" err: ", fmt.Errorf("wrong password"))
		return
	}

	tx, err := parseTransaction(cmd, args)
	if err != nil {
		log.Errorf(cmd.Use+" err: %s", err.Error())
		return
	}
	account, err := AccountByRpc(tx.From().String())
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	if tx.Header.Nonce == 0 {
		tx.Header.Nonce = account.Nonce + 1
	}
	if err := signMsg(tx, privKey.Private); err != nil {
		log.Error(cmd.Use+" err: ", errors.New("signature failure"))
		return
	}

	rs, err := sendMsg(tx)
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
	} else if rs.Code != 0 {
		log.Errorf(cmd.Use+" err: code %d, message: %s", rs.Code, rs.Err)
	} else {
		fmt.Println()
		fmt.Println(string(rs.Result))
	}
}

func parseTransaction(cmd *cobra.Command, args []string) (*types.Message, error) {
	var err error
	var from, to, token arry.Address
	var amount, fee, nonce uint64
	from = arry.StringToAddress(args[0])
	to = arry.StringToAddress(args[1])
	token = arry.StringToAddress(args[2])
	if fAmount, err := strconv.ParseFloat(args[3], 64); err != nil {
		return nil, errors.New("[amount] wrong")
	} else {
		if fAmount < 0 {
			return nil, errors.New("[amount] wrong")
		}
		if amount, err = amount2.NewAmount(fAmount); err != nil {
			return nil, errors.New("[amount] wrong")
		}
	}
	if fFees, err := strconv.ParseFloat(args[4], 64); err != nil {
		return nil, errors.New("[fees] wrong")
	} else {
		if fFees < 0 {
			return nil, errors.New("[fees] wrong")
		}
		if fee, err = amount2.NewAmount(fFees); err != nil {
			return nil, errors.New("[fees] wrong")
		}
	}
	if len(args) > 6 {
		nonce, err = strconv.ParseUint(args[6], 10, 64)
		if err != nil {
			return nil, errors.New("[nonce] wrong")
		}
	}
	return message.NewTransaction(from, to, token, amount, fee, nonce), nil
}

func signMsg(msg *types.Message, key string) error {
	msg.SetHash()
	priv, err := secp256k1.ParseStringToPrivate(key)
	if err != nil {
		return errors.New("[key] wrong")
	}
	if err := msg.SignMessage(priv); err != nil {
		return errors.New("sign failed")
	}
	return nil
}

func sendMsg(msg *types.Message) (*rpc.Response, error) {
	rpcMsg, err := rpctypes.MsgToRpcMsg(msg)
	if err != nil {
		return nil, err
	}
	rpcClient, err := NewRpcClient()
	if err != nil {
		return nil, err
	}
	defer rpcClient.Close()

	jsonBytes, err := json.Marshal(rpcMsg)
	if err != nil {
		return nil, err
	}
	re := &rpc.SendMessageCode{Code: jsonBytes}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()

	resp, err := rpcClient.Gc.SendMessageRaw(ctx, re)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

var GetMessageCmd = &cobra.Command{
	Use:     "GetMessage {msghash}; Get Message by hash;",
	Aliases: []string{"getmessage", "GM", "gm"},
	Short:   "GetMessage {msghash}; Get Message by hash;",
	Example: `
	GetMessage 0xef7b92e552dca02c97c9d596d1bf69d0044d95dec4cee0e6a20153e62bce893b
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  GetMessage,
}

func GetMessage(cmd *cobra.Command, args []string) {
	resp, err := GetMessageRpc(args[0])
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	if resp.Code == 0 {
		output(string(resp.Result))
		return
	}
	outputRespError(cmd.Use, resp)
}

func GetMessageRpc(hashStr string) (*rpc.Response, error) {
	client, err := NewRpcClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()


	resp, err := client.Gc.GetMessage(ctx, &rpc.Hash{Hash: hashStr})
	return resp, err
}
