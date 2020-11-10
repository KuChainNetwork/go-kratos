package types

import (
	"github.com/KuChainNetwork/kuchain/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgBlock struct {
	FullBlock
}

type MsgTx struct {
	sdk.TxResponse
}

type MsgMsg struct {
	types.KuMsg
}
