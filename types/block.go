package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

type TxData sdk.TxResponse

type FullBlock struct {
	Header     `json:"header"`
	BlockID    string               `json:"block_id"`
	Txs        []TxData             `json:"txs"`
	Evidence   tmTypes.EvidenceData `json:"evidence"`
	LastCommit tmTypes.Commit       `json:"last_commit"`
}
