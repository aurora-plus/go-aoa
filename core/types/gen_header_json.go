package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/Aurorachain/go-Aurora/common"
	"github.com/Aurorachain/go-Aurora/common/hexutil"
)

var _ = (*headerMarshaling)(nil)

func (h Header) MarshalJSON() ([]byte, error) {
	type Header struct {
		ParentHash         common.Hash    `json:"parentHash"       gencodec:"required"`
		Coinbase           common.Address `json:"miner"            gencodec:"required"`
		Root               common.Hash    `json:"stateRoot"        gencodec:"required"`
		TxHash             common.Hash    `json:"transactionsRoot" gencodec:"required"`
		ReceiptHash        common.Hash    `json:"receiptsRoot"     gencodec:"required"`
		Bloom              Bloom          `json:"logsBloom"        gencodec:"required"`
		Number             *hexutil.Big   `json:"number"           gencodec:"required"`
		GasLimit           hexutil.Uint64 `json:"gasLimit"         gencodec:"required"`
		GasUsed            hexutil.Uint64 `json:"gasUsed"          gencodec:"required"`
		Time               *hexutil.Big   `json:"timestamp"        gencodec:"required"`
		Extra              hexutil.Bytes  `json:"extraData"        gencodec:"required"`
		AgentName          []byte         `json:"agentName"        gencodec:"required"`
		DelegateRoot       common.Hash    `json:"delegateRoot"     gencodec:"required"`
		ShuffleHash        common.Hash    `json:"shuffleHash"      gencodec:"required"`
		ShuffleBlockNumber *big.Int       `json:"shuffleBlockNumber"        gencodec:"required"`
		Hash               common.Hash    `json:"hash"`
	}
	var enc Header
	enc.ParentHash = h.ParentHash
	enc.Coinbase = h.Coinbase
	enc.Root = h.Root
	enc.TxHash = h.TxHash
	enc.ReceiptHash = h.ReceiptHash
	enc.Bloom = h.Bloom
	enc.Number = (*hexutil.Big)(h.Number)
	enc.GasLimit = hexutil.Uint64(h.GasLimit)
	enc.GasUsed = hexutil.Uint64(h.GasUsed)
	enc.Time = (*hexutil.Big)(h.Time)
	enc.Extra = h.Extra
	enc.AgentName = h.AgentName
	enc.DelegateRoot = h.DelegateRoot
	enc.ShuffleHash = h.ShuffleHash
	enc.ShuffleBlockNumber = h.ShuffleBlockNumber
	enc.Hash = h.Hash()
	return json.Marshal(&enc)
}

func (h *Header) UnmarshalJSON(input []byte) error {
	type Header struct {
		ParentHash         *common.Hash    `json:"parentHash"       gencodec:"required"`
		Coinbase           *common.Address `json:"miner"            gencodec:"required"`
		Root               *common.Hash    `json:"stateRoot"        gencodec:"required"`
		TxHash             *common.Hash    `json:"transactionsRoot" gencodec:"required"`
		ReceiptHash        *common.Hash    `json:"receiptsRoot"     gencodec:"required"`
		Bloom              *Bloom          `json:"logsBloom"        gencodec:"required"`
		Number             *hexutil.Big    `json:"number"           gencodec:"required"`
		GasLimit           *hexutil.Uint64 `json:"gasLimit"         gencodec:"required"`
		GasUsed            *hexutil.Uint64 `json:"gasUsed"          gencodec:"required"`
		Time               *hexutil.Big    `json:"timestamp"        gencodec:"required"`
		Extra              *hexutil.Bytes  `json:"extraData"        gencodec:"required"`
		AgentName          []byte          `json:"agentName"        gencodec:"required"`
		DelegateRoot       *common.Hash    `json:"delegateRoot"     gencodec:"required"`
		ShuffleHash        *common.Hash    `json:"shuffleHash"      gencodec:"required"`
		ShuffleBlockNumber *big.Int        `json:"shuffleBlockNumber"        gencodec:"required"`
	}
	var dec Header
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.ParentHash == nil {
		return errors.New("missing required field 'parentHash' for Header")
	}
	h.ParentHash = *dec.ParentHash
	if dec.Coinbase == nil {
		return errors.New("missing required field 'miner' for Header")
	}
	h.Coinbase = *dec.Coinbase
	if dec.Root == nil {
		return errors.New("missing required field 'stateRoot' for Header")
	}
	h.Root = *dec.Root
	if dec.TxHash == nil {
		return errors.New("missing required field 'transactionsRoot' for Header")
	}
	h.TxHash = *dec.TxHash
	if dec.ReceiptHash == nil {
		return errors.New("missing required field 'receiptsRoot' for Header")
	}
	h.ReceiptHash = *dec.ReceiptHash
	if dec.Bloom == nil {
		return errors.New("missing required field 'logsBloom' for Header")
	}
	h.Bloom = *dec.Bloom
	if dec.Number == nil {
		return errors.New("missing required field 'number' for Header")
	}
	h.Number = (*big.Int)(dec.Number)
	if dec.GasLimit == nil {
		return errors.New("missing required field 'gasLimit' for Header")
	}
	h.GasLimit = uint64(*dec.GasLimit)
	if dec.GasUsed == nil {
		return errors.New("missing required field 'gasUsed' for Header")
	}
	h.GasUsed = uint64(*dec.GasUsed)
	if dec.Time == nil {
		return errors.New("missing required field 'timestamp' for Header")
	}
	h.Time = (*big.Int)(dec.Time)
	if dec.Extra == nil {
		return errors.New("missing required field 'extraData' for Header")
	}
	h.Extra = *dec.Extra
	if dec.AgentName == nil {
		return errors.New("missing required field 'agentName' for Header")
	}
	h.AgentName = dec.AgentName
	if dec.DelegateRoot == nil {
		return errors.New("missing required field 'delegateRoot' for Header")
	}
	h.DelegateRoot = *dec.DelegateRoot
	if dec.ShuffleHash == nil {
		return errors.New("missing required field 'shuffleHash' for Header")
	}
	h.ShuffleHash = *dec.ShuffleHash
	if dec.ShuffleBlockNumber == nil {
		return errors.New("missing required field 'shuffleBlockNumber' for Header")
	}
	h.ShuffleBlockNumber = dec.ShuffleBlockNumber
	return nil
}
