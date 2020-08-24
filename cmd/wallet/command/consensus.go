package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/Futuremine-chain/futuremine/futuremine/common/kit/message"
	private2 "github.com/Futuremine-chain/futuremine/futuremine/common/private"
	"github.com/Futuremine-chain/futuremine/futuremine/rpc"
	"github.com/Futuremine-chain/futuremine/futuremine/types"
	"github.com/Futuremine-chain/futuremine/service/p2p"
	"github.com/Futuremine-chain/futuremine/tools/amount"
	"github.com/Futuremine-chain/futuremine/tools/crypto/ecc/secp256k1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

func init() {
	txCmds := []*cobra.Command{
		SendCandidateCmd,
		SendCancelCmd,
		SendVoteCmd,
		GetCandidatesCmd,
		CycleSupersCmd,
	}
	RootCmd.AddCommand(txCmds...)
	RootSubCmdGroups["consensus"] = txCmds

}

var SendCandidateCmd = &cobra.Command{
	Use:     "SendCandidate {address} {fees} {password} {nonce}; Become candidate;",
	Aliases: []string{"sendcandidate", "SC", "sc"},
	Short:   "SendCandidate {address} {fees} {password} {nonce}; Become candidate;",
	Example: `
	SendCandidate xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ 0.001
		OR
	SendCandidate xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ 0.001 123456
		OR
	SendCandidate xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ 0.001 123456 1
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
	if err := signMsg(candidateMsg, privKey.String()); err != nil {
		log.Error(cmd.Use+" err: ", errors.New("signature failure"))
		return
	}

	rs, err := sendMsg(candidateMsg)
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
	var from string
	var fee, nonce uint64
	from = args[0]

	if fFees, err := strconv.ParseFloat(args[1], 64); err != nil {
		return nil, errors.New("[fees] wrong")
	} else {
		if fFees < 0 {
			return nil, errors.New("[fees] wrong")
		}
		if fee, err = amount.NewAmount(fFees); err != nil {
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
	SendCancel xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ 0.001
		OR
	SendCancel xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ 0.001 123456
		OR
	SendCancel xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ 0.001 123456 1
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

	cancel, err := parseCancel(args)
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
	if err := signMsg(cancel, privKey.Private); err != nil {
		log.Error(cmd.Use+" err: ", errors.New("signature failure"))
		return
	}

	rs, err := sendMsg(cancel)
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
	} else if rs.Code != 0 {
		log.Errorf(cmd.Use+" err: code %d, message: %s", rs.Code, rs.Err)
	} else {
		fmt.Println()
		fmt.Println(string(rs.Result))
	}
}

func parseCancel(args []string) (*types.Message, error) {
	var err error
	var from string
	var fee, nonce uint64
	from = args[0]
	if fFees, err := strconv.ParseFloat(args[1], 64); err != nil {
		return nil, errors.New("[fees] wrong")
	} else {
		if fFees < 0 {
			return nil, errors.New("[fees] wrong")
		}
		if fee, err = amount.NewAmount(fFees); err != nil {
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
	SendVote xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ xCE9boXz2TxSE9srVPDdfszyiXtfT3vduc8 0.001
		OR
	SendVote xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ xCE9boXz2TxSE9srVPDdfszyiXtfT3vduc8 0.001 123456
		OR
	SendVote xCHiGPLCzgnrdTqjKABXZteAGVJu3jXLjnQ xCE9boXz2TxSE9srVPDdfszyiXtfT3vduc8 0.001 123456 1
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
	if err := signMsg(vote, privKey.Private); err != nil {
		log.Error(cmd.Use+" err: ", errors.New("signature failure"))
		return
	}

	rs, err := sendMsg(vote)
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
	var from, to string
	var fee, nonce uint64
	from = args[0]
	to = args[1]
	if fFees, err := strconv.ParseFloat(args[2], 64); err != nil {
		return nil, errors.New("[fees] wrong")
	} else {
		if fFees < 0 {
			return nil, errors.New("[fees] wrong")
		}
		if fee, err = amount.NewAmount(fFees); err != nil {
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

var GetCandidatesCmd = &cobra.Command{
	Use:     "GetCandidates",
	Short:   "GetCandidates;Get current candidates;",
	Aliases: []string{"getcandidates", "GC", "gc"},
	Example: `
	GetCandidates
	`,
	Run: GetCandidates,
}

func GetCandidates(cmd *cobra.Command, args []string) {
	client, err := NewRpcClient()
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()
	resp, err := client.Gc.Candidates(ctx, &rpc.NullReq{})
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

var CycleSupersCmd = &cobra.Command{
	Use:     "CycleSupers {cycle}; Gets the current super nodes;",
	Short:   "CycleSupers {cycle}; Gets the current super nodes;",
	Aliases: []string{"cyclesupers", "CS", "cs"},
	Example: `
	CycleSupers {8736163}
	`,
	Args: cobra.MinimumNArgs(1),
	Run:  CycleSupers,
}

func CycleSupers(cmd *cobra.Command, args []string) {
	term, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		log.Error(cmd.Use+" err: ", errors.New("[term] wrong"))
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

	resp, err := client.Gc.GetCycleSupers(ctx, &rpc.CycleReq{Cycle: term})
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
