package core

import (
	"encoding/json"
	"errors"

	"github.com/Aurorachain/go-Aurora/common"
	"github.com/Aurorachain/go-Aurora/common/hexutil"
	"github.com/Aurorachain/go-Aurora/common/math"
	"github.com/Aurorachain/go-Aurora/params"
)

var _ = (*genesisSpecMarshaling)(nil)

func (g Genesis) MarshalJSON() ([]byte, error) {
	type Genesis struct {
		Config     *params.ChainConfig                         `json:"config"`
		Timestamp  math.HexOrDecimal64                         `json:"timestamp"`
		ExtraData  hexutil.Bytes                               `json:"extraData"`
		GasLimit   math.HexOrDecimal64                         `json:"gasLimit"   gencodec:"required"`
		Coinbase   common.Address                              `json:"coinbase"`
		Alloc      map[common.UnprefixedAddress]GenesisAccount `json:"alloc"      gencodec:"required"`
		Agents     GenesisAgents                               `json:"agents"     gencodec:"required"`
		Number     math.HexOrDecimal64                         `json:"number"`
		GasUsed    math.HexOrDecimal64                         `json:"gasUsed"`
		ParentHash common.Hash                                 `json:"parentHash"`
	}
	var enc Genesis
	enc.Config = g.Config
	enc.Timestamp = math.HexOrDecimal64(g.Timestamp)
	enc.ExtraData = g.ExtraData
	enc.GasLimit = math.HexOrDecimal64(g.GasLimit)
	enc.Coinbase = g.Coinbase
	if g.Alloc != nil {
		enc.Alloc = make(map[common.UnprefixedAddress]GenesisAccount, len(g.Alloc))
		for k, v := range g.Alloc {
			enc.Alloc[common.UnprefixedAddress(k)] = v
		}
	}
	enc.Agents = g.Agents
	enc.Number = math.HexOrDecimal64(g.Number)
	enc.GasUsed = math.HexOrDecimal64(g.GasUsed)
	enc.ParentHash = g.ParentHash
	return json.Marshal(&enc)
}

func (g *Genesis) UnmarshalJSON(input []byte) error {
	type Genesis struct {
		Config     *params.ChainConfig                         `json:"config"`
		Timestamp  *math.HexOrDecimal64                        `json:"timestamp"`
		ExtraData  *hexutil.Bytes                              `json:"extraData"`
		GasLimit   *math.HexOrDecimal64                        `json:"gasLimit"   gencodec:"required"`
		Coinbase   *common.Address                             `json:"coinbase"`
		Alloc      map[common.UnprefixedAddress]GenesisAccount `json:"alloc"      gencodec:"required"`
		Agents     *GenesisAgents                              `json:"agents"     gencodec:"required"`
		Number     *math.HexOrDecimal64                        `json:"number"`
		GasUsed    *math.HexOrDecimal64                        `json:"gasUsed"`
		ParentHash *common.Hash                                `json:"parentHash"`
	}
	var dec Genesis
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Config != nil {
		g.Config = dec.Config
	}
	if dec.Timestamp != nil {
		g.Timestamp = uint64(*dec.Timestamp)
	}
	if dec.ExtraData != nil {
		g.ExtraData = *dec.ExtraData
	}
	if dec.GasLimit == nil {
		return errors.New("missing required field 'gasLimit' for Genesis")
	}
	g.GasLimit = uint64(*dec.GasLimit)
	if dec.Coinbase != nil {
		g.Coinbase = *dec.Coinbase
	}
	if dec.Alloc == nil {
		return errors.New("missing required field 'alloc' for Genesis")
	}
	g.Alloc = make(GenesisAlloc, len(dec.Alloc))
	for k, v := range dec.Alloc {
		g.Alloc[common.Address(k)] = v
	}
	if dec.Agents == nil {
		return errors.New("missing required field 'agents' for Genesis")
	}
	g.Agents = *dec.Agents
	if dec.Number != nil {
		g.Number = uint64(*dec.Number)
	}
	if dec.GasUsed != nil {
		g.GasUsed = uint64(*dec.GasUsed)
	}
	if dec.ParentHash != nil {
		g.ParentHash = *dec.ParentHash
	}
	return nil
}