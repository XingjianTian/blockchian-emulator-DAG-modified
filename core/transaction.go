// Definition of transaction

package core

import (
	"blockEmulator/utils"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"time"
)

/*
交易在 BlockEmulator 中，包含了 7 个基本的字段，包含交易的发起者、交易的接受者、交易的 Nonce 值、签名、交易的值、交易的哈希以及交易加入到交易池的时间（用于统计交易的接受时间）。

同时因为在 BlockEmulator 中预设了不少跨分片相关的设计，需要在交易中增加特别的功能，如 Relayed 字段表示是否被 relay 池所接收，这也是 relay 跨分片交易机制所需要的功能，如字段 OriginalSender 、 字段FinalRecipient 以及字段 RawTxHash 都是用于 BrokerChain 跨分片交易机制所需要的特别设计
*/

type Transaction struct {
	Sender    utils.Address
	Recipient utils.Address
	Nonce     uint64
	Signature []byte // not implemented now.
	Value     *big.Int
	TxHash    []byte

	Time time.Time // the time adding in pool

	// used in transaction relaying
	Relayed bool
	// used in broker, if the tx is not a broker1 or broker2 tx, these values should be empty.
	HasBroker      bool
	SenderIsBroker bool
	OriginalSender utils.Address
	FinalRecipient utils.Address
	RawTxHash      []byte
}

func (tx *Transaction) PrintTx() string {
	vals := []interface{}{
		tx.Sender[:],
		tx.Recipient[:],
		tx.Value,
		string(tx.TxHash[:]),
	}
	res := fmt.Sprintf("%v\n", vals)
	return res
}

// Encode transaction for storing
func (tx *Transaction) Encode() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// Decode transaction
func DecodeTx(to_decode []byte) *Transaction {
	var tx Transaction

	decoder := gob.NewDecoder(bytes.NewReader(to_decode))
	err := decoder.Decode(&tx)
	if err != nil {
		log.Panic(err)
	}

	return &tx
}

// new a transaction
func NewTransaction(sender, recipient string, value *big.Int, nonce uint64) *Transaction {
	tx := &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Value:     value,
		Nonce:     nonce,
	}

	hash := sha256.Sum256(tx.Encode())
	tx.TxHash = hash[:]
	tx.Relayed = false
	tx.FinalRecipient = ""
	tx.OriginalSender = ""
	tx.RawTxHash = nil
	tx.HasBroker = false
	tx.SenderIsBroker = false
	return tx
}
