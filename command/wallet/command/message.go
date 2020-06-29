package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Futuremine-chain/futuremine/futuremine/common/kit/message"
	private2 "github.com/Futuremine-chain/futuremine/futuremine/common/private"
	"github.com/Futuremine-chain/futuremine/futuremine/rpc"
	rpctypes "github.com/Futuremine-chain/futuremine/futuremine/rpc/types"
	"github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/service/p2p"
	"github.com/Futuremine-chain/futuremine/tools/arry"
	"github.com/Futuremine-chain/futuremine/tools/crypto/ecc/secp256k1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

func init() {
	txCmds := []*cobra.Command{
		GetMessageCmd,
		SendMessageCmd,
		SendCandidateCmd,
		SendCancelCmd,
		SendVoteCmd,
	}
	RootCmd.AddCommand(txCmds...)
	RootSubCmdGroups["message"] = txCmds

}

var SendMessageCmd = &cobra.Command{
	Use:     "SendTransction {from} {to} {token} {amount} {fees} {password} {nonce}; Send a transaction;",
	Aliases: []string{"sendtransction", "ST", "st"},
	Short:   "SendTransction {from} {to} {token} {amount} {fees} {password} {nonce}; Send a transaction;",
	Example: `
	SendTransaction 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ UB 10 0.1
		OR
	SendTransaction 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ UB 10 0.1 123456
		OR
	SendTransaction 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ UB 10 0.1 123456 1
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
	if !signMsg(cmd, tx, privKey.Private) {
		log.Error(cmd.Use+" err: ", errors.New("signature failure"))
		return
	}

	rs, err := sendMsg(cmd, tx)
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
		if amount, err = types.NewAmount(fAmount); err != nil {
			return nil, errors.New("[amount] wrong")
		}
	}
	if fFees, err := strconv.ParseFloat(args[4], 64); err != nil {
		return nil, errors.New("[fees] wrong")
	} else {
		if fFees < 0 {
			return nil, errors.New("[fees] wrong")
		}
		if fee, err = types.NewAmount(fFees); err != nil {
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

func signMsg(cmd *cobra.Command, tx *types.Message, key string) bool {
	priv, err := secp256k1.ParseStringToPrivate(key)
	if err != nil {
		log.Error(cmd.Use+" err: ", errors.New("[key] wrong"))
		return false
	}
	if err := tx.SignMessage(priv); err != nil {
		log.Error(cmd.Use+" err: ", errors.New("sign failed"))
		return false
	}
	return true
}

func sendMsg(cmd *cobra.Command, msg *types.Message) (*rpc.Response, error) {
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
	re := &rpc.Request{Params: jsonBytes}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()

	resp, err := rpcClient.Gc.SendMessageRaw(ctx, re)
	if err != nil {
		return nil, err
	}
	return resp, nil

}

var SendCandidateCmd = &cobra.Command{
	Use:     "SendCandidate {address} {fees} {password} {nonce}; Become candidate;",
	Aliases: []string{"sendcandidate", "SC", "sc"},
	Short:   "SendCandidate {address} {fees} {password} {nonce}; Become candidate;",
	Example: `
	SendCandidate 3ajF4MdbBYE2UPESEyhQbdUj2Y28CNwGDCWA 0.001
		OR
	SendCandidate 3ajF4MdbBYE2UPESEyhQbdUj2Y28CNwGDCWA 0.001 123456
		OR
	SendCandidate 3ajF4MdbBYE2UPESEyhQbdUj2Y28CNwGDCWA 0.001 123456 1
`,
	Args: cobra.MinimumNArgs(2),
	Run:  SendCandidate,
}

func SendCandidate(cmd *cobra.Command, args []string) {
	var passwd []byte
	var err error
	if len(args) > 2 {
		passwd = []byte(args[2])
	} else {
		fmt.Println("please input password：")
		passwd, err = readPassWd()
		if err != nil {
			log.Error(cmd.Use+" err: ", fmt.Errorf("read password failed! %s", err.Error()))
			return
		}
	}
	private, err := loadPrivate(getAddJsonPath(args[0]), passwd)
	if err != nil {
		log.Error(cmd.Use+" err: ", fmt.Errorf("wrong password"))
		return
	}
	privKey, err := secp256k1.ParseStringToPrivate(private.Private)
	if err != nil {
		log.Error(cmd.Use+" err: ", fmt.Errorf("failed to parse private %s", err.Error()))
		return
	}
	p2pId, _ := p2p.PrivateToP2pId(private2.NewPrivate(privKey))

	candidateMsg, err := parseCandidate(cmd, args, p2pId.String())
	if err != nil {
		log.Error(cmd.Use+" err: ", err.Error())
		return
	}
	account, err := AccountByRpc(candidateMsg.From().String())
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	if candidateMsg.Header.Nonce == 0 {
		candidateMsg.Header.Nonce = account.Nonce + 1
	}
	if !signMsg(cmd, candidateMsg, privKey.String()) {
		log.Error(cmd.Use+" err: ", errors.New("signature failure"))
		return
	}

	rs, err := sendMsg(cmd, candidateMsg)
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
	} else if rs.Code != 0 {
		log.Errorf(cmd.Use+" err: code %d, message: %s", rs.Code, rs.Err)
	} else {
		fmt.Println()
		fmt.Println(string(rs.Result))
	}
}

func parseCandidate(cmd *cobra.Command, args []string, p2pid string) (*types.Message, error) {
	var err error
	var from arry.Address
	var fee, nonce uint64
	from = arry.StringToAddress(args[0])

	if fFees, err := strconv.ParseFloat(args[1], 64); err != nil {
		return nil, errors.New("[fees] wrong")
	} else {
		if fFees < 0 {
			return nil, errors.New("[fees] wrong")
		}
		if fee, err = types.NewAmount(fFees); err != nil {
			log.Error(cmd.Use + " err: ")
			return nil, errors.New("[fees] wrong")
		}
	}
	if len(args) > 3 {
		nonce, err = strconv.ParseUint(args[3], 10, 64)
		if err != nil {
			return nil, errors.New("[nonce] wrong")
		}
	}

	return message.NewCandidate(from, p2pid, fee, nonce), nil
}

var SendCancelCmd = &cobra.Command{
	Use:     "SendCancel {address} {fees} {password} {nonce}; Cancel candidate;",
	Aliases: []string{"sendcancel", "SCL", "scl"},
	Short:   "SendCancel {address} {fees} {password} {nonce}; Cancel candidate;",
	Example: `
	SendCancel 3ajF4MdbBYE2UPESEyhQbdUj2Y28CNwGDCWA 0.001
		OR
	SendCancel 3ajF4MdbBYE2UPESEyhQbdUj2Y28CNwGDCWA 0.001 123456
		OR
	SendCancel 3ajF4MdbBYE2UPESEyhQbdUj2Y28CNwGDCWA 0.001 123456 1
	`,
	Args: cobra.MinimumNArgs(2),
	Run:  CancelCandidate,
}

func CancelCandidate(cmd *cobra.Command, args []string) {
	var passwd []byte
	var err error
	if len(args) > 2 {
		passwd = []byte(args[2])
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

		return
	}

	cancel, err := parseCancel(cmd, args)
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	account, err := AccountByRpc(cancel.From().String())
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	if cancel.Header.Nonce == 0 {
		cancel.Header.Nonce = account.Nonce + 1
	}
	if !signMsg(cmd, cancel, privKey.Private) {
		log.Error(cmd.Use+" err: ", errors.New("signature failure"))
		return
	}

	rs, err := sendMsg(cmd, cancel)
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
	} else if rs.Code != 0 {
		log.Errorf(cmd.Use+" err: code %d, message: %s", rs.Code, rs.Err)
	} else {
		fmt.Println()
		fmt.Println(string(rs.Result))
	}
}

func parseCancel(cmd *cobra.Command, args []string) (*types.Message, error) {
	var err error
	var from arry.Address
	var fee, nonce uint64
	from = arry.StringToAddress(args[0])
	if fFees, err := strconv.ParseFloat(args[1], 64); err != nil {
		return nil, errors.New("[fees] wrong")
	} else {
		if fFees < 0 {
			return nil, errors.New("[fees] wrong")
		}
		if fee, err = types.NewAmount(fFees); err != nil {
			return nil, errors.New("[fees] wrong")
		}
	}
	if len(args) > 3 {
		nonce, err = strconv.ParseUint(args[3], 10, 64)
		if err != nil {
			return nil, errors.New("[nonce] wrong")
		}
	}
	return message.NewCancel(from, fee, nonce), nil
}

var SendVoteCmd = &cobra.Command{
	Use:     "SendVote {from} {to} {fees} {password} {nonce}；Vote for a candidate;",
	Aliases: []string{"sendvote", "SV", "sv"},
	Short:   "SendVote {from} {to} {fees} {password} {nonce}; Vote for a candidate;",
	Example: `
	Vote 3ajF4MdbBYE2UPESEyhQbdUj2Y28CNwGDCWA 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 0.001
		OR
	Vote 3ajF4MdbBYE2UPESEyhQbdUj2Y28CNwGDCWA 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 0.001 123456
		OR
	Vote 3ajF4MdbBYE2UPESEyhQbdUj2Y28CNwGDCWA 3ajDJUnMYDyzXLwefRfNp7yLcdmg3ULb9ndQ 0.001 123456 1
`,
	Args: cobra.MinimumNArgs(3),
	Run:  Vote,
}

func Vote(cmd *cobra.Command, args []string) {
	var passwd []byte
	var err error
	if len(args) > 3 {
		passwd = []byte(args[3])
	} else {
		fmt.Println("please input password：")
		passwd, err = readPassWd()
		if err != nil {

			return
		}
	}
	privKey, err := loadPrivate(getAddJsonPath(args[0]), passwd)
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}

	vote, err := parseVote(args)
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	account, err := AccountByRpc(vote.From().String())
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}

	if vote.Header.Nonce == 0 {
		vote.Header.Nonce = account.Nonce + 1
	}
	if !signMsg(cmd, vote, privKey.Private) {
		log.Error(cmd.Use+" err: ", errors.New("signature failure"))
		return
	}

	rs, err := sendMsg(cmd, vote)
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
	} else if rs.Code != 0 {
		log.Errorf(cmd.Use+" err: code %d, message: %s", rs.Code, rs.Err)
	} else {
		fmt.Println()
		fmt.Println(string(rs.Result))
	}
}

func parseVote(args []string) (*types.Message, error) {
	var err error
	var from, to arry.Address
	var fee, nonce uint64
	from = arry.StringToAddress(args[0])
	to = arry.StringToAddress(args[1])
	if fFees, err := strconv.ParseFloat(args[2], 64); err != nil {
		return nil, errors.New("[fees] wrong")
	} else {
		if fFees < 0 {
			return nil, errors.New("[fees] wrong")
		}
		if fee, err = types.NewAmount(fFees); err != nil {
			return nil, errors.New("[fees] wrong")
		}
	}
	if len(args) > 4 {
		nonce, err = strconv.ParseUint(args[4], 10, 64)
		if err != nil {
			return nil, errors.New("[nonce] wrong")
		}
	}
	return message.NewVote(from, to, fee, nonce), nil
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

	params := []interface{}{hashStr}
	if bytes, err := json.Marshal(params); err != nil {
		return nil, err
	} else {
		resp, err := client.Gc.GetMessage(ctx, &rpc.Request{Params: bytes})
		return resp, err
	}
}
