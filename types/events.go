package types

import tmTypes "github.com/tendermint/tendermint/types"

const (
	EventNewBlockStr = "tendermint/event/NewBlock"
)

type EventNewBlock struct {
	Type  string                    `json:"type"`
	Block tmTypes.EventDataNewBlock `json:"value"`
}
