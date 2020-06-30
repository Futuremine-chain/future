package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Futuremine-chain/futuremine/futuremine/common/kit"
	"github.com/Futuremine-chain/futuremine/futuremine/common/kit/message"
	"github.com/Futuremine-chain/futuremine/futuremine/rpc"
	"github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

func init() {
	contractCmds := []*cobra.Command{
		TokenCmd,
		SendCreateTokenCmd,
	}
	RootCmd.AddCommand(contractCmds...)
	RootSubCmdGroups["token"] = contractCmds

}

var SendCreateTokenCmd = &cobra.Command{
	Use:     "SendCreateToken {from} {to} {name} {shorthand} {amount} {fees} {password} {nonce}; Send and create token;",
	Aliases: []string{"sendcreatetoken", "sct"},
	Short:   "SendCreateToken {from} {to} {name} {shorthand} {amount} {fees} {password} {nonce}; Send and create token;",
	Example: `
	SendCreateToken 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE "M token" MT 1000 0.1
		OR
	SendCreateToken 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE "M token" MT 1000 0.1 123456
		OR
	SendCreateToken 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajNkh7yVYkETL9JKvGx3aL2YVNrqksjCUUE "M token" MT 1000 0.1 123456 0
	`,
	Args: cobra.MinimumNArgs(6),
	Run:  SendCreateToken,
}

func SendCreateToken(cmd *cobra.Command, args []string) {
	var passwd []byte
	var err error
	if len(args) > 6 {
		passwd = []byte(args[6])
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

	tokenMsg, err := parseToken(args)
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	account, err := AccountByRpc(tokenMsg.From().String())
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	if tokenMsg.Header.Nonce == 0 {
		tokenMsg.Header.Nonce = account.Nonce + 1
	}
	if err := signMsg(tokenMsg, privKey.Private); err != nil {
		log.Error(cmd.Use+" err: ", errors.New("signature failure"))
		return
	}

	rs, err := sendMsg(tokenMsg)
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
	} else if rs.Code != 0 {
		log.Errorf(cmd.Use+" err: code %d, message: %s", rs.Code, rs.Err)
	} else {
		fmt.Println()
		fmt.Println(string(rs.Result))
	}
}

func parseToken(args []string) (*types.Message, error) {
	var err error
	var from, to, tokenAddr arry.Address
	var amount, fee, nonce uint64
	var name, shorthand string
	from = arry.StringToAddress(args[0])
	to = arry.StringToAddress(args[1])
	name = args[2]
	shorthand = args[3]
	if fAmount, err := strconv.ParseFloat(args[4], 64); err != nil {
		return nil, errors.New("[amount] wrong")
	} else {
		if fAmount < 0 {
			return nil, errors.New("[amount] wrong")
		}
		if amount, err = types.NewAmount(fAmount); err != nil {
			return nil, errors.New("[amount] wrong")
		}
	}
	tokenAddr, err = kit.GenerateTokenAddress(Net, from, shorthand)
	if err != nil {
		return nil, err
	}
	fmt.Println("token address is ", tokenAddr.String())

	if fFees, err := strconv.ParseFloat(args[5], 64); err != nil {
		return nil, errors.New("[fees] wrong")
	} else {
		if fFees < 0 {
			return nil, errors.New("[fees] wrong")
		}
		if fee, err = types.NewAmount(fFees); err != nil {
			return nil, errors.New("[fees] wrong")
		}
	}
	if len(args) > 7 {
		nonce, err = strconv.ParseUint(args[7], 10, 64)
		if err != nil {
			return nil, errors.New("[nonce] wrong")
		}
	}
	tokenMsg := message.NewToken(from, to, tokenAddr, amount, fee, nonce, name, shorthand)
	return tokenMsg, nil
}

var TokenCmd = &cobra.Command{
	Use:     "Token {token address}; Get a token;",
	Aliases: []string{"token", "T", "t"},
	Short:   "Token {token address}; Get a token;",
	Example: `
	GetContract RtH8MnY2yVxHPFDSkV9YeQ5igrppSNtnmVp
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  Token,
}

func Token(cmd *cobra.Command, args []string) {
	resp, err := GetTokenByRpc(args[0])
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	if resp.Code == 0 {
		output(string(resp.Result))
		return
	} else {
		outputRespError(cmd.Use, resp)
	}
}

func GetTokenByRpc(tokenAddr string) (*rpc.Response, error) {
	client, err := NewRpcClient()
	if err != nil {
		return nil, err
	}
	defer client.Close()

	params := []interface{}{tokenAddr}
	if bytes, err := json.Marshal(params); err != nil {
		return nil, err
	} else {
		re := &rpc.Request{Params: bytes}
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
		defer cancel()
		return client.Gc.Token(ctx, re)
	}
}
