package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/Aurorachain/go-Aurora/common"
	"github.com/Aurorachain/go-Aurora/common/hexutil"
)

var _ = (*txdataMarshaling)(nil)

func (t txdata) MarshalJSON() ([]byte, error) {
	type txdata struct {
		AccountNonce hexutil.Uint64  `json:"nonce"    gencodec:"required"`
		Price        *hexutil.Big    `json:"gasPrice" gencodec:"required"`
		GasLimit     hexutil.Uint64  `json:"gas"      gencodec:"required"`
		Recipient    *common.Address `json:"to"       rlp:"nil"`
		Amount       *hexutil.Big    `json:"value"    gencodec:"required"`
		Payload      hexutil.Bytes   `json:"input"    gencodec:"required"`
		Action       uint64          `json:"action"  gencodec:"required"`
		Vote         []byte          `json:"vote" rlp:"nil"`
		Nickname     []byte          `json:"nickname" rlp:"nil"`
		Asset        *common.Address `json:"asset,omitempty" rlp:"nil"`
		AssetInfo    []byte          `json:"assetInfo,omitempty" rlp:"nil"`
		SubAddress   string          `json:"subAddress,omitempty" rlp:"nil"`
		Abi          string          `json:"abi,omitempty" rlp:"nil"`
		V            *hexutil.Big    `json:"v" gencodec:"required"`
		R            *hexutil.Big    `json:"r" gencodec:"required"`
		S            *hexutil.Big    `json:"s" gencodec:"required"`
		Hash         *common.Hash    `json:"hash" rlp:"-"`
	}
	var enc txdata
	enc.AccountNonce = hexutil.Uint64(t.AccountNonce)
	enc.Price = (*hexutil.Big)(t.Price)
	enc.GasLimit = hexutil.Uint64(t.GasLimit)
	enc.Recipient = t.Recipient
	enc.Amount = (*hexutil.Big)(t.Amount)
	enc.Payload = t.Payload
	enc.Action = t.Action
	enc.Vote = t.Vote
	enc.Nickname = t.Nickname
	enc.Asset = t.Asset
	enc.AssetInfo = t.AssetInfo
	enc.SubAddress = t.SubAddress
	enc.Abi = t.Abi
	enc.V = (*hexutil.Big)(t.V)
	enc.R = (*hexutil.Big)(t.R)
	enc.S = (*hexutil.Big)(t.S)
	enc.Hash = t.Hash
	return json.Marshal(&enc)
}

func (t *txdata) UnmarshalJSON(input []byte) error {
	type txdata struct {
		AccountNonce *hexutil.Uint64 `json:"nonce"    gencodec:"required"`
		Price        *hexutil.Big    `json:"gasPrice" gencodec:"required"`
		GasLimit     *hexutil.Uint64 `json:"gas"      gencodec:"required"`
		Recipient    *common.Address `json:"to"       rlp:"nil"`
		Amount       *hexutil.Big    `json:"value"    gencodec:"required"`
		Payload      *hexutil.Bytes  `json:"input"    gencodec:"required"`
		Action       *uint64         `json:"action"  gencodec:"required"`
		Vote         []byte          `json:"vote" rlp:"nil"`
		Nickname     []byte          `json:"nickname" rlp:"nil"`
		Asset        *common.Address `json:"asset,omitempty" rlp:"nil"`
		AssetInfo    []byte          `json:"assetInfo,omitempty" rlp:"nil"`
		SubAddress   *string         `json:"subAddress,omitempty" rlp:"nil"`
		Abi          *string         `json:"abi,omitempty" rlp:"nil"`
		V            *hexutil.Big    `json:"v" gencodec:"required"`
		R            *hexutil.Big    `json:"r" gencodec:"required"`
		S            *hexutil.Big    `json:"s" gencodec:"required"`
		Hash         *common.Hash    `json:"hash" rlp:"-"`
	}
	var dec txdata
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.AccountNonce == nil {
		return errors.New("missing required field 'nonce' for txdata")
	}
	t.AccountNonce = uint64(*dec.AccountNonce)
	if dec.Price == nil {
		return errors.New("missing required field 'gasPrice' for txdata")
	}
	t.Price = (*big.Int)(dec.Price)
	if dec.GasLimit == nil {
		return errors.New("missing required field 'gas' for txdata")
	}
	t.GasLimit = uint64(*dec.GasLimit)
	if dec.Recipient != nil {
		t.Recipient = dec.Recipient
	}
	if dec.Amount == nil {
		return errors.New("missing required field 'value' for txdata")
	}
	t.Amount = (*big.Int)(dec.Amount)
	if dec.Payload == nil {
		return errors.New("missing required field 'input' for txdata")
	}
	t.Payload = *dec.Payload
	if dec.Action == nil {
		return errors.New("missing required field 'action' for txdata")
	}
	t.Action = *dec.Action
	if dec.Vote != nil {
		t.Vote = dec.Vote
	}
	if dec.Nickname != nil {
		t.Nickname = dec.Nickname
	}
	if dec.Asset != nil {
		t.Asset = dec.Asset
	}
	if dec.AssetInfo != nil {
		t.AssetInfo = dec.AssetInfo
	}
	if dec.SubAddress != nil {
		t.SubAddress = *dec.SubAddress
	}
	if dec.Abi != nil {
		t.Abi = *dec.Abi
	}
	if dec.V == nil {
		return errors.New("missing required field 'v' for txdata")
	}
	t.V = (*big.Int)(dec.V)
	if dec.R == nil {
		return errors.New("missing required field 'r' for txdata")
	}
	t.R = (*big.Int)(dec.R)
	if dec.S == nil {
		return errors.New("missing required field 's' for txdata")
	}
	t.S = (*big.Int)(dec.S)
	if dec.Hash != nil {
		t.Hash = dec.Hash
	}
	return nil
}
