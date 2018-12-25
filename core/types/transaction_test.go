package types

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"testing"

	"encoding/binary"
	"fmt"
	"github.com/Aurorachain/go-Aurora/common"
	"github.com/Aurorachain/go-Aurora/common/hexutil"
	"github.com/Aurorachain/go-Aurora/crypto"
	"github.com/Aurorachain/go-Aurora/rlp"
	"io"
	"sync/atomic"
	"regexp"
)

var (
	emptyTx = NewTransaction(
		0,
		common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		big.NewInt(0), 0, big.NewInt(0),
		nil, 0, nil, nil, nil, nil, "")

)

func TestTransactionSigHash(t *testing.T) {

}

func TestTransactionEncode(t *testing.T) {

}

func decodeTx(data []byte) (*Transaction, error) {
	var tx Transaction
	t, err := &tx, rlp.Decode(bytes.NewReader(data), &tx)

	return t, err
}

func defaultTestKey() (*ecdsa.PrivateKey, common.Address) {
	key, _ := crypto.HexToECDSA("45a915e4d060149eb4365960e6a7a45f334393093061116b197e3240065ff2d8")
	addr := crypto.PubkeyToAddress(key.PublicKey)
	return key, addr
}

func TestRecipientEmpty(t *testing.T) {
	_, addr := defaultTestKey()
	tx, err := decodeTx(common.Hex2Bytes("f8498080808080011ca09b16de9d5bdee2cf56c28d16275a4da68cd30273e2525f3959f5d62557489921a0372ebd8fb3345f7db7b5a86d42e24d36e983e259b0664ceb8c227ec9af572f3d"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	from, err := Sender(AuroraSigner{}, tx)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if addr != from {
		t.Error("derived address doesn't match")
	}
}

func TestRecipientNormal(t *testing.T) {
	_, addr := defaultTestKey()

	tx, err := decodeTx(common.Hex2Bytes("f85d80808094000000000000000000000000000000000000000080011ca0527c0d8f5c63f7b9f41324a7c8a563ee1190bcbf0dac8ab446291bdbf32f5c79a0552c4ef0a09a04395074dab9ed34d3fbfb843c2f2546cc30fe89ec143ca94ca6"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	from, err := Sender(AuroraSigner{}, tx)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if addr != from {
		t.Error("derived address doesn't match")
	}
}

func TestTransactionPriceNonceSort(t *testing.T) {

	keys := make([]*ecdsa.PrivateKey, 25)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}

	signer := AuroraSigner{}

	groups := map[common.Address]Transactions{}
	for start, key := range keys {
		addr := crypto.PubkeyToAddress(key.PublicKey)
		for i := 0; i < 25; i++ {
			tx, _ := SignTx(NewTransaction(uint64(start+i), common.Address{}, big.NewInt(100), 100, big.NewInt(int64(start+i)), nil, 0, nil, nil, nil, nil, ""), signer, key)
			groups[addr] = append(groups[addr], tx)
		}
	}

	txset := NewTransactionsByPriceAndNonce(signer, groups)

	txs := Transactions{}
	for tx := txset.Peek(); tx != nil; tx = txset.Peek() {
		txs = append(txs, tx)
		txset.Shift()
	}
	if len(txs) != 25*25 {
		t.Errorf("expected %d transactions, found %d", 25*25, len(txs))
	}
	for i, txi := range txs {
		fromi, _ := Sender(signer, txi)

		for j, txj := range txs[i+1:] {
			fromj, _ := Sender(signer, txj)

			if fromi == fromj && txi.Nonce() > txj.Nonce() {
				t.Errorf("invalid nonce ordering: tx #%d (A=%x N=%v) < tx #%d (A=%x N=%v)", i, fromi[:4], txi.Nonce(), i+j, fromj[:4], txj.Nonce())
			}
		}

		prev, next := i-1, i+1
		for j := i - 1; j >= 0; j-- {
			if fromj, _ := Sender(signer, txs[j]); fromi == fromj {
				prev = j
				break
			}
		}
		for j := i + 1; j < len(txs); j++ {
			if fromj, _ := Sender(signer, txs[j]); fromi == fromj {
				next = j
				break
			}
		}

		for j := prev + 1; j < next; j++ {
			fromj, _ := Sender(signer, txs[j])
			if j < i && txs[j].GasPrice().Cmp(txi.GasPrice()) < 0 {
				t.Errorf("invalid gasprice ordering: tx #%d (A=%x P=%v) < tx #%d (A=%x P=%v)", j, fromj[:4], txs[j].GasPrice(), i, fromi[:4], txi.GasPrice())
			}
			if j > i && txs[j].GasPrice().Cmp(txi.GasPrice()) > 0 {
				t.Errorf("invalid gasprice ordering: tx #%d (A=%x P=%v) > tx #%d (A=%x P=%v)", j, fromj[:4], txs[j].GasPrice(), i, fromi[:4], txi.GasPrice())
			}
		}
	}
}

func TestTransactionJSON(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("could not generate key: %v", err)
	}

	signer := AuroraSigner{chainId: common.Big1}

	for i := uint64(0); i < 25; i++ {
		var tx *Transaction
		switch i % 2 {
		case 0:
			tx = NewTransaction(i, common.Address{1}, common.Big0, 1, common.Big2, []byte("abcdef"), 0, nil, nil, nil, nil, "")
		case 1:

			tx = NewContractCreation(i, common.Big0, 1, common.Big2, []byte("abcdef"), `[{"constant":true,"inputs":[],"name":"mybalance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"payable":true,"stateMutability":"payable","type":"fallback"}]`,nil)
		}

		tx, err := SignTx(tx, signer, key)
		if err != nil {
			t.Fatalf("could not sign transaction: %v", err)
		}

		data, err := json.Marshal(tx)
		if err != nil {
			t.Errorf("json.Marshal failed: %v", err)
		}

		var parsedTx *Transaction
		if err := json.Unmarshal(data, &parsedTx); err != nil {
			t.Errorf("json.Unmarshal failed: %v", err)
		}

		if tx.Hash() != parsedTx.Hash() {
			t.Errorf("parsed tx differs from original tx, want %v, got %v", tx, parsedTx)
		}
		if tx.ChainId().Cmp(parsedTx.ChainId()) != 0 {
			t.Errorf("invalid chain id, want %d, got %d", tx.ChainId(), parsedTx.ChainId())
		}
	}
}

func TestTransaction_DecodeRLP(t *testing.T) {

	errSign := "0xf872808502540be40083030d40943e106d2004a5bdc48be21c28e46c9e0c2d28d69f8ad3c21bcecceda1000000801ca0b2df725d4f5647ea4199d375e0fc1bb17363ee551a3f96842cb4817b5e35b57ca055c8ab5061c3865768838ac8f5734071327dce29042a595a4088df5379a973f2808080"

	encodedTx, _ := hexutil.Decode(errSign)
	fmt.Println(encodedTx)
	tx := new(transaction)

	if err := rlp.DecodeBytes(encodedTx, tx); err != nil {
		fmt.Printf("rlp decode error:%v trx:%v\n", err, tx)
	} else {
		fmt.Println(tx)
	}
}

func TestTxDifference(t *testing.T) {
	v1 := Vote{nil, 1}
	json.Marshal(v1)
	v2 := Vote{nil, 2}
	v3 := Vote{nil, 3}

	buf := new(bytes.Buffer)
	var data = []interface{}{
		v1,
		v2,
		v3,
	}
	for _, v := range data {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
	}
	fmt.Printf("%x", buf.Bytes())
	fmt.Println(buf.Bytes())
}

type transaction struct {
	data tdata

	hash atomic.Value
	size atomic.Value
	from atomic.Value
}

type tdata struct {
	AccountNonce uint64          `json:"nonce"    gencodec:"required"`
	Price        *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit     uint64          `json:"gas"      gencodec:"required"`
	Recipient    *common.Address `json:"to"       rlp:"nil"`
	Amount       *big.Int        `json:"value"    gencodec:"required"`
	Payload      []byte          `json:"input"    gencodec:"required"`

	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	Action   uint   `json:"action"  gencodec:"required"`
	Vote     []byte `json:"vote" rlp:"nil"`
	Nickname []byte `json:"nickname" rlp:"nil"`

	Hash *common.Hash `json:"hash" rlp:"-"`

	AssetSymbol string `json:"assetSymbol,omitempty" rlp:"nil"`

	AssetInfo *AssetInfo `json:"assetInfo,omitempty" rlp:"nil"`
}

func (tx *transaction) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &tx.data)
}

func (tx *transaction) DecodeRLP(s *rlp.Stream) error {
	_, size, _ := s.Kind()
	err := s.Decode(&tx.data)
	if err == nil {
		tx.size.Store(common.StorageSize(rlp.ListSize(size)))
	}

	return err
}

func TestDeriveSha(t *testing.T) {
	testString := []string{"12", "23", "23"}
	fmt.Println(testString)
	enc, err := json.Marshal(testString)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(enc)
	var dec []string
	err = json.Unmarshal(enc, dec)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(dec)
}

func TestDecode(t *testing.T) {
	enc, err := hexutil.Decode("0xf8b10784ee6b2800830186a0940dff1151bb88110b679babcfd4599656203df7cc80b8441d834a1b0000000000000000000000000000000000000000000000000000000000000032000000000000000000000000000000000000000000000000000000000000003206808080808080819ba0a00cbcfe542320ed91cbfafb52f9a420e4909110d7aaf80d53c163f0d1f58e78a034b5057cc30bdd8a48afcd2890c7bacd91e2d76aefee2284492390027307578e")
	if err != nil {
		fmt.Println(err)
		return
	}
	tx := new(Transaction)
	err = rlp.DecodeBytes(enc, tx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tx)
}

func TestRecoverFromAddress(t *testing.T) {
	enc, err := hexutil.Decode("0xf8738084ee6b280083015f90945d17e0168071c56882b403451966a31cb508d6c9880b1a2bc2ec500000808080808080808025a030abd95402a9df641009e6a8a805df9795fc8704e2b576eaf0b6361eae4389aba0146818e069e3c3fdad2c798a2e3c9569e45bcf3608345bbbf7bbe8018ca53b81")
	if err != nil {
		t.Fatal(err)
	}
	tx := new(Transaction)
	err = rlp.DecodeBytes(enc, tx)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tx)

}

func TestAssetInfoToBytes(t *testing.T) {
	var assetInfo AssetInfo
	assetInfo.Desc = "aa"
	assetInfo.Issuer = &common.Address{}
	assetInfo.Name = "AOA"
	assetInfo.Symbol = "A"
	assetInfo.Supply = big.NewInt(1111)

	toBytes, err := AssetInfoToBytes(assetInfo)
	if err != nil {
		t.Fatal(err)
	}

	assetInfo2, err := BytesToAssetInfo(toBytes)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(assetInfo2)

}

func TestAoaAddress(t *testing.T) {
	to := "AOA60aac5adbb14ea09b3a01f04b56aa8b5db420f55"
	match, err := regexp.MatchString("(?i:^AOA|0x)[0-9a-f]{40}[0-9A-Za-z]{0,32}$", to)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(match)
}
