// the define and some operation of txpool

package core

import (
	"blockEmulator/utils"
	"sync"
	"time"
)

type TxPool struct {
	TxQueue   []*Transaction            // transaction Queue
	RelayPool map[uint64][]*Transaction //designed for sharded blockchain, from Monoxide
	lock      sync.Mutex
	// The pending list is ignored
	excutedMap map[string]bool
}

func NewTxPool() *TxPool {
	return &TxPool{
		TxQueue:    make([]*Transaction, 0),
		RelayPool:  make(map[uint64][]*Transaction),
		excutedMap: make(map[string]bool),
	}
}

// Add a transaction to the pool (consider the queue only)
func (txpool *TxPool) AddTx2Pool(tx *Transaction) {
	txpool.lock.Lock()
	defer txpool.lock.Unlock()
	if tx.Time.IsZero() {
		tx.Time = time.Now()
	}
	txpool.TxQueue = append(txpool.TxQueue, tx)
}

// Add a list of transactions to the pool
func (txpool *TxPool) AddTxs2Pool(txs []*Transaction) {
	txpool.lock.Lock()
	defer txpool.lock.Unlock()

	var txsToAdd = InterSect(txs, txpool.TxQueue)

	for _, tx := range txsToAdd {
		if tx.Time.IsZero() {
			tx.Time = time.Now()
		}
		txpool.TxQueue = append(txpool.TxQueue, tx)
	}
}

// add transactions into the pool head
func (txpool *TxPool) AddTxs2Pool_Head(tx []*Transaction) {
	txpool.lock.Lock()
	defer txpool.lock.Unlock()
	txpool.TxQueue = append(tx, txpool.TxQueue...)
}

// Pack transactions for a proposal
func (txpool *TxPool) PackTxs(max_txs uint64) []*Transaction {
	txpool.lock.Lock()
	defer txpool.lock.Unlock()
	txNum := max_txs
	if uint64(len(txpool.TxQueue)) < txNum {
		txNum = uint64(len(txpool.TxQueue))
	}
	txs_Packed := txpool.TxQueue[:txNum]
	txpool.TxQueue = txpool.TxQueue[txNum:]
	return txs_Packed
}

// 返回 slice1 中存在，但在 slice2 中不存在的元素
func InterSect(slice1, slice2 []*Transaction) []*Transaction {
	// 创建一个用于存储结果的map
	m := make(map[string]bool)

	// 将slice2中的元素设置为map的key，值无所谓
	for _, item := range slice2 {
		m[string(item.TxHash)] = true
	}

	// 遍历slice1，如果元素不在map中，则加入结果切片
	var diff []*Transaction
	for _, item := range slice1 {
		if !m[string(item.TxHash)] {
			diff = append(diff, item)
		}
	}

	return diff
}

// Pack transactions for a proposal
func (txpool *TxPool) RemoveTxs(excutedTxs []*Transaction) {
	txpool.lock.Lock()
	defer txpool.lock.Unlock()

	txpool.TxQueue = InterSect(txpool.TxQueue, excutedTxs)
}

// Relay transactions
func (txpool *TxPool) AddRelayTx(tx *Transaction, shardID uint64) {
	txpool.lock.Lock()
	defer txpool.lock.Unlock()
	_, ok := txpool.RelayPool[shardID]
	if !ok {
		txpool.RelayPool[shardID] = make([]*Transaction, 0)
	}
	txpool.RelayPool[shardID] = append(txpool.RelayPool[shardID], tx)
}

// txpool get locked
func (txpool *TxPool) GetLocked() {
	txpool.lock.Lock()
}

// txpool get unlocked
func (txpool *TxPool) GetUnlocked() {
	txpool.lock.Unlock()
}

// get the length of tx queue
func (txpool *TxPool) GetTxQueueLen() int {
	txpool.lock.Lock()
	defer txpool.lock.Unlock()
	return len(txpool.TxQueue)
}

// get the length of ClearRelayPool
func (txpool *TxPool) ClearRelayPool() {
	txpool.lock.Lock()
	defer txpool.lock.Unlock()
	txpool.RelayPool = nil
}

// abort ! Pack relay transactions from relay pool
func (txpool *TxPool) PackRelayTxs(shardID, minRelaySize, maxRelaySize uint64) ([]*Transaction, bool) {
	txpool.lock.Lock()
	defer txpool.lock.Unlock()
	if _, ok := txpool.RelayPool[shardID]; !ok {
		return nil, false
	}
	if len(txpool.RelayPool[shardID]) < int(minRelaySize) {
		return nil, false
	}
	txNum := maxRelaySize
	if uint64(len(txpool.RelayPool[shardID])) < txNum {
		txNum = uint64(len(txpool.RelayPool[shardID]))
	}
	relayTxPacked := txpool.RelayPool[shardID][:txNum]
	txpool.RelayPool[shardID] = txpool.RelayPool[shardID][txNum:]
	return relayTxPacked, true
}

// abort ! Transfer transactions when re-sharding
func (txpool *TxPool) TransferTxs(addr utils.Address) []*Transaction {
	txpool.lock.Lock()
	defer txpool.lock.Unlock()
	txTransfered := make([]*Transaction, 0)
	newTxQueue := make([]*Transaction, 0)
	for _, tx := range txpool.TxQueue {
		if tx.Sender == addr {
			txTransfered = append(txTransfered, tx)
		} else {
			newTxQueue = append(newTxQueue, tx)
		}
	}
	newRelayPool := make(map[uint64][]*Transaction)
	for shardID, shardPool := range txpool.RelayPool {
		for _, tx := range shardPool {
			if tx.Sender == addr {
				txTransfered = append(txTransfered, tx)
			} else {
				if _, ok := newRelayPool[shardID]; !ok {
					newRelayPool[shardID] = make([]*Transaction, 0)
				}
				newRelayPool[shardID] = append(newRelayPool[shardID], tx)
			}
		}
	}
	txpool.TxQueue = newTxQueue
	txpool.RelayPool = newRelayPool
	return txTransfered
}
