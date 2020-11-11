package sdk

import (
	"fmt"

	"github.com/KuChainNetwork/go-kratos/types"
	"github.com/KuChainNetwork/kuchain/chain/client/rest/block"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

func (c *Client) QueryBlockByNum(num int64) (block.DecodeResultBlock, error) {
	var res block.DecodeResultBlock

	if err := c.queryFromJSON(&res, c.LcdURL, "blocks/%d/decode", num); err != nil {
		return res, errors.Wrapf(err, "error by query block decode %d", num)
	}

	return res, nil
}

func (c *Client) QueryLatestBlock() (block.DecodeResultBlock, error) {
	var res block.DecodeResultBlock

	if err := c.queryFromJSON(&res, c.LcdURL, "blocks/latest/decode"); err != nil {
		return res, errors.Wrapf(err, "error by query block decode latest")
	}

	return res, nil
}

func (c *Client) QueryTxByHash(hash string) (sdk.TxResponse, error) {
	var tx struct {
		Res sdk.TxResponse `json:"TxResponse"`
	}

	if err := c.queryFromJSON(&tx, c.LcdURL, "txs/%s", hash); err != nil {
		return tx.Res, errors.Wrapf(err, "error by query tx %s", hash)
	}

	return tx.Res, nil
}

func (c *Client) QueryFullBlock(num int64) (types.FullBlock, error) {
	var decodeRes block.DecodeResultBlock

	path := fmt.Sprintf("blocks/%d/decode", num)
	if num == 0 {
		path = "blocks/latest/decode"
	}

	if err := c.queryFromJSON(&decodeRes, c.LcdURL, path); err != nil {
		return types.FullBlock{}, errors.Wrapf(err, "error by query block decode %d", num)
	}

	res := types.FullBlock{
		Header:     decodeRes.DecodeBlock.Header,
		BlockID:    decodeRes.BlockID.Hash.String(),
		Evidence:   decodeRes.DecodeBlock.Evidence,
		LastCommit: *decodeRes.DecodeBlock.LastCommit,
		TxDatas:    make([]types.TxData, 0, len(decodeRes.DecodeBlock.TxsHash)),
		Txs:        decodeRes.DecodeBlock.Txs,
		TxsHash:    decodeRes.DecodeBlock.TxsHash,
	}

	// get each tx
	for _, txHash := range decodeRes.DecodeBlock.TxsHash {
		txRes, err := c.QueryTxByHash(txHash.String())
		if err != nil {
			return types.FullBlock{}, errors.Wrapf(err, "error by tx query %s", txHash.String())
		}

		res.TxDatas = append(res.TxDatas, types.TxData(txRes))
	}

	return res, nil
}
