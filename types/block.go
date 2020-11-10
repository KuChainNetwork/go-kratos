package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmTypes "github.com/tendermint/tendermint/types"
)

type TxData sdk.TxResponse

type FullBlock struct {
	Header     `json:"header"`
	BlockID    string             `json:"block_id"`
	TxDatas    []TxData           `json:"tx_datas"`
	Txs        []json.RawMessage  `json:"txs"`
	TxsHash    []tmbytes.HexBytes `json:"txs_hash"`
	hash       tmbytes.HexBytes
	Evidence   tmTypes.EvidenceData `json:"evidence"`
	LastCommit tmTypes.Commit       `json:"last_commit"`
}
