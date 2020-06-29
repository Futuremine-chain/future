package command

import (
	"context"
	"encoding/json"
	"github.com/Futuremine-chain/futuremine/futuremine/rpc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strconv"
	"time"
)

func init() {
	blockCmds := []*cobra.Command{GetBlockCmd}

	RootCmd.AddCommand(blockCmds...)
	RootSubCmdGroups["block"] = blockCmds
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

	params := []interface{}{height}
	if bytes, err := json.Marshal(params); err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	} else {
		resp, err := client.Gc.GetBlockHeight(ctx, &rpc.Request{Params: bytes})
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

	params := []interface{}{args[0]}
	if bytes, err := json.Marshal(params); err != nil {
		log.Error(cmd.Use+" err: ", err)
		return
	} else {
		resp, err := client.Gc.GetBlockHash(ctx, &rpc.Request{Params: bytes})
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
