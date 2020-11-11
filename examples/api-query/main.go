package main

import (
	"context"
	"flag"
	"fmt"

	sdk "github.com/KuChainNetwork/go-kratos"
	chainCfg "github.com/KuChainNetwork/kuchain/chain/config"
	"github.com/KuChainNetwork/kuchain/utils/log"
)

var (
	num = flag.Int64("num", 0, "block height to query")
	url = flag.String("url", "http://127.0.0.1:1317/", "kuchain lcd url")
	rpc = flag.String("rpc", "http://127.0.0.1:26657/", "kuchain rpc url")
)

func main() {
	flag.Parse()

	chainCfg.SealChainConfig()

	ctx := sdk.NewCtx(context.Background()).
		WithUrls(*url, *rpc).
		WithLogger(log.NewLoggerByZap(true, "*:debug"))

	cli := sdk.NewClient(ctx)

	data, err := cli.QueryFullBlock(*num)

	if err != nil {
		panic(fmt.Errorf("data err by %s", err.Error()))
	}

	fmt.Println(string(sdk.MustMarshalJSON(data)))

	blockResults, err := cli.QueryBlockResultsByHeight(*num)
	if err != nil {
		panic(fmt.Errorf("result err by %s", err.Error()))
	}

	fmt.Println(string(sdk.MustMarshalJSON(blockResults)))
}
