package command

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Futuremine-chain/futuremine/futuremine/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

func init() {
	nodeCmds := []*cobra.Command{
		LastHeightCmd,
		MsgPoolCmd,
		PeerInfoCmd,
		GetCandidatesCmd,
		CycleSupersCmd,
		LocalInfoCmd,
	}
	RootCmd.AddCommand(nodeCmds...)
	RootSubCmdGroups["node"] = nodeCmds
}

//GenerateCmd cpu mine block
var MsgPoolCmd = &cobra.Command{
	Use:     "MsgPool",
	Short:   "MsgPool; Get messages in the message pool;",
	Aliases: []string{"msgpool", "MP", "mp"},
	Example: `
	MsgPool 
	`,
	Args: cobra.MinimumNArgs(0),
	Run:  MsgPool,
}

func MsgPool(cmd *cobra.Command, args []string) {

	client, err := NewRpcClient()
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()

	resp, err := client.Gc.GetMsgPool(ctx, &rpc.Request{Params: nil})
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
	resp, err := client.Gc.Candidates(ctx, &rpc.Request{})
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

	params := []interface{}{term}
	if bytes, err := json.Marshal(params); err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	} else {
		resp, err := client.Gc.GetCycleSupers(ctx, &rpc.Request{Params: bytes})
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
}

var PeerInfoCmd = &cobra.Command{
	Use:     "PeerInfo",
	Short:   "PeerInfo; Get peer info;",
	Aliases: []string{"peerinfo", "PI", "pi"},
	Example: `
	PeerInfo 
	`,
	Args: cobra.MinimumNArgs(0),
	Run:  PeerInfo,
}

func PeerInfo(cmd *cobra.Command, args []string) {
	client, err := NewRpcClient()
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()
	resp, err := client.Gc.PeersInfo(ctx, &rpc.Request{})
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
	resp, err := client.Gc.LastHeight(ctx, &rpc.Request{})
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

var LocalInfoCmd = &cobra.Command{
	Use:     "LocalInfo ;Get the current node information",
	Short:   "LocalInfo ;Get the current node information;",
	Aliases: []string{"localinfo", "LI", "li"},
	Example: `
	LocalInfo
	`,
	Args: cobra.MinimumNArgs(0),
	Run:  LocalInfo,
}

func LocalInfo(cmd *cobra.Command, args []string) {
	client, err := NewRpcClient()
	if err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*20)
	defer cancel()
	resp, err := client.Gc.LocalInfo(ctx, &rpc.Request{})
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
