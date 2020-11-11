package types

import (
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	sdk "github.com/tendermint/tendermint/types"
)

type (
	Header   = sdk.Header
	HexBytes = tmbytes.HexBytes
)

type (
	ResultBlockResults = ctypes.ResultBlockResults
)
