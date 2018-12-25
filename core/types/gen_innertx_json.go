package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/Aurorachain/go-Aurora/common"
	"github.com/Aurorachain/go-Aurora/common/hexutil"
)

var _ = (*innertxMarshaling)(nil)

func (i InnerTx) MarshalJSON() ([]byte, error) {
	type InnerTx struct {
		From    common.Address  `json:"from" gencodec:"required"`
		To      common.Address  `json:"to" gencodec:"to" gencodec:"required"`
		AssetID *common.Address `json:"assetid" rlp:"nil"`
		Value   *hexutil.Big    `json:"value" gencodec:"required"`
	}
	var enc InnerTx
	enc.From = i.From
	enc.To = i.To
	enc.AssetID = i.AssetID
	enc.Value = (*hexutil.Big)(i.Value)
	return json.Marshal(&enc)
}

func (i *InnerTx) UnmarshalJSON(input []byte) error {
	type InnerTx struct {
		From    *common.Address `json:"from" gencodec:"required"`
		To      *common.Address `json:"to" gencodec:"to" gencodec:"required"`
		AssetID *common.Address `json:"assetid" rlp:"nil"`
		Value   *hexutil.Big    `json:"value" gencodec:"required"`
	}
	var dec InnerTx
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.From == nil {
		return errors.New("missing required field 'from' for InnerTx")
	}
	i.From = *dec.From
	if dec.To != nil {
		i.To = *dec.To
	}
	if dec.AssetID != nil {
		i.AssetID = dec.AssetID
	}
	if dec.Value == nil {
		return errors.New("missing required field 'value' for InnerTx")
	}
	i.Value = (*big.Int)(dec.Value)
	return nil
}
