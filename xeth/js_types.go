package xeth

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethutil"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/state"
)

// Block interface exposed to QML
type JSBlock struct {
	//Transactions string `json:"transactions"`
	ref          *types.Block
	Size         string        `json:"size"`
	Number       int           `json:"number"`
	Hash         string        `json:"hash"`
	Transactions *ethutil.List `json:"transactions"`
	Uncles       *ethutil.List `json:"uncles"`
	Time         int64         `json:"time"`
	Coinbase     string        `json:"coinbase"`
	Name         string        `json:"name"`
	GasLimit     string        `json:"gasLimit"`
	GasUsed      string        `json:"gasUsed"`
	PrevHash     string        `json:"prevHash"`
	Bloom        string        `json:"bloom"`
	Raw          string        `json:"raw"`
}

// Creates a new QML Block from a chain block
func NewJSBlock(block *types.Block) *JSBlock {
	if block == nil {
		return &JSBlock{}
	}

	ptxs := make([]*JSTransaction, len(block.Transactions()))
	for i, tx := range block.Transactions() {
		ptxs[i] = NewJSTx(tx)
	}
	txlist := ethutil.NewList(ptxs)

	puncles := make([]*JSBlock, len(block.Uncles()))
	for i, uncle := range block.Uncles() {
		puncles[i] = NewJSBlock(types.NewBlockWithHeader(uncle))
	}
	ulist := ethutil.NewList(puncles)

	return &JSBlock{
		ref: block, Size: block.Size().String(),
		Number: int(block.NumberU64()), GasUsed: block.GasUsed().String(),
		GasLimit: block.GasLimit().String(), Hash: ethutil.Bytes2Hex(block.Hash()),
		Transactions: txlist, Uncles: ulist,
		Time:     block.Time(),
		Coinbase: ethutil.Bytes2Hex(block.Coinbase()),
		PrevHash: ethutil.Bytes2Hex(block.ParentHash()),
		Bloom:    ethutil.Bytes2Hex(block.Bloom()),
		Raw:      block.String(),
	}
}

func (self *JSBlock) ToString() string {
	if self.ref != nil {
		return self.ref.String()
	}

	return ""
}

func (self *JSBlock) GetTransaction(hash string) *JSTransaction {
	tx := self.ref.Transaction(ethutil.Hex2Bytes(hash))
	if tx == nil {
		return nil
	}

	return NewJSTx(tx)
}

type JSTransaction struct {
	ref *types.Transaction

	Value           string `json:"value"`
	Gas             string `json:"gas"`
	GasPrice        string `json:"gasPrice"`
	Hash            string `json:"hash"`
	Address         string `json:"address"`
	Sender          string `json:"sender"`
	RawData         string `json:"rawData"`
	Data            string `json:"data"`
	Contract        bool   `json:"isContract"`
	CreatesContract bool   `json:"createsContract"`
	Confirmations   int    `json:"confirmations"`
}

func NewJSTx(tx *types.Transaction) *JSTransaction {
	hash := ethutil.Bytes2Hex(tx.Hash())
	receiver := ethutil.Bytes2Hex(tx.To())
	if receiver == "0000000000000000000000000000000000000000" {
		receiver = ethutil.Bytes2Hex(core.AddressFromMessage(tx))
	}
	sender := ethutil.Bytes2Hex(tx.From())
	createsContract := core.MessageCreatesContract(tx)

	var data string
	if createsContract {
		data = strings.Join(core.Disassemble(tx.Data()), "\n")
	} else {
		data = ethutil.Bytes2Hex(tx.Data())
	}

	return &JSTransaction{ref: tx, Hash: hash, Value: ethutil.CurrencyToString(tx.Value()), Address: receiver, Contract: createsContract, Gas: tx.Gas().String(), GasPrice: tx.GasPrice().String(), Data: data, Sender: sender, CreatesContract: createsContract, RawData: ethutil.Bytes2Hex(tx.Data())}
}

func (self *JSTransaction) ToString() string {
	return self.ref.String()
}

type JSKey struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

func NewJSKey(key *crypto.KeyPair) *JSKey {
	return &JSKey{ethutil.Bytes2Hex(key.Address()), ethutil.Bytes2Hex(key.PrivateKey), ethutil.Bytes2Hex(key.PublicKey)}
}

type JSObject struct {
	*Object
}

func NewJSObject(object *Object) *JSObject {
	return &JSObject{object}
}

type PReceipt struct {
	CreatedContract bool   `json:"createdContract"`
	Address         string `json:"address"`
	Hash            string `json:"hash"`
	Sender          string `json:"sender"`
}

func NewPReciept(contractCreation bool, creationAddress, hash, address []byte) *PReceipt {
	return &PReceipt{
		contractCreation,
		ethutil.Bytes2Hex(creationAddress),
		ethutil.Bytes2Hex(hash),
		ethutil.Bytes2Hex(address),
	}
}

// Peer interface exposed to QML

type JSPeer struct {
	ref     *p2p.Peer
	Ip      string `json:"ip"`
	Version string `json:"version"`
	Caps    string `json:"caps"`
}

func NewJSPeer(peer *p2p.Peer) *JSPeer {
	var caps []string
	for _, cap := range peer.Caps() {
		caps = append(caps, fmt.Sprintf("%s/%d", cap.Name, cap.Version))
	}

	return &JSPeer{
		ref:     peer,
		Ip:      fmt.Sprintf("%v", peer.RemoteAddr()),
		Version: fmt.Sprintf("%v", peer.Identity()),
		Caps:    fmt.Sprintf("%v", caps),
	}
}

type JSReceipt struct {
	CreatedContract bool   `json:"createdContract"`
	Address         string `json:"address"`
	Hash            string `json:"hash"`
	Sender          string `json:"sender"`
}

func NewJSReciept(contractCreation bool, creationAddress, hash, address []byte) *JSReceipt {
	return &JSReceipt{
		contractCreation,
		ethutil.Bytes2Hex(creationAddress),
		ethutil.Bytes2Hex(hash),
		ethutil.Bytes2Hex(address),
	}
}

type JSMessage struct {
	To        string `json:"to"`
	From      string `json:"from"`
	Input     string `json:"input"`
	Output    string `json:"output"`
	Path      int32  `json:"path"`
	Origin    string `json:"origin"`
	Timestamp int32  `json:"timestamp"`
	Coinbase  string `json:"coinbase"`
	Block     string `json:"block"`
	Number    int32  `json:"number"`
	Value     string `json:"value"`
}

func NewJSMessage(message *state.Message) JSMessage {
	return JSMessage{
		To:        ethutil.Bytes2Hex(message.To),
		From:      ethutil.Bytes2Hex(message.From),
		Input:     ethutil.Bytes2Hex(message.Input),
		Output:    ethutil.Bytes2Hex(message.Output),
		Path:      int32(message.Path),
		Origin:    ethutil.Bytes2Hex(message.Origin),
		Timestamp: int32(message.Timestamp),
		Coinbase:  ethutil.Bytes2Hex(message.Origin),
		Block:     ethutil.Bytes2Hex(message.Block),
		Number:    int32(message.Number.Int64()),
		Value:     message.Value.String(),
	}
}
