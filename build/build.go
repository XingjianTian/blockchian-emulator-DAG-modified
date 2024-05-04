package build

import (
	"blockEmulator/consensus_shard/pbft_all"
	"blockEmulator/core"
	"blockEmulator/params"
	"blockEmulator/supervisor"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"
)

func initConfig(nid, nnm, sid, snm uint64) *params.ChainConfig {
	params.ShardNum = int(snm)
	for i := uint64(0); i < snm; i++ {
		if _, ok := params.IPmap_nodeTable[i]; !ok {
			params.IPmap_nodeTable[i] = make(map[uint64]string)
		}
		for j := uint64(0); j < nnm; j++ {
			params.IPmap_nodeTable[i][j] = "127.0.0.1:" + strconv.Itoa(28800+int(i)*100+int(j))
		}
	}
	params.IPmap_nodeTable[params.DeciderShard] = make(map[uint64]string)
	params.IPmap_nodeTable[params.DeciderShard][0] = params.SupervisorAddr
	params.NodesInShard = int(nnm)
	params.ShardNum = int(snm)

	pcc := &params.ChainConfig{
		ChainID:        sid,
		NodeID:         nid,
		ShardID:        sid,
		Nodes_perShard: uint64(params.NodesInShard),
		ShardNums:      snm,
		BlockSize:      uint64(params.MaxBlockSize_global),
		BlockInterval:  uint64(params.Block_Interval),
		InjectSpeed:    uint64(params.InjectSpeed),
	}
	return pcc
}

func BuildSupervisor(nnm, snm, mod uint64, isDAG bool) {
	var measureMod []string
	if mod == 0 || mod == 2 {
		measureMod = params.MeasureBrokerMod
	} else {
		measureMod = params.MeasureRelayMod
	}

	lsn := new(supervisor.Supervisor)
	lsn.NewSupervisor(isDAG, params.SupervisorAddr, initConfig(123, nnm, 123, snm), params.CommitteeMethod[mod], measureMod...)
	time.Sleep(10000 * time.Millisecond)
	go lsn.SupervisorTxHandling()
	lsn.TcpListen()
}

func BuildNewPbftNode(nid, nnm, sid, snm, mod uint64, isDAG bool) {
	worker := pbft_all.NewPbftNode(sid, nid, initConfig(nid, nnm, sid, snm), params.CommitteeMethod[mod], isDAG)

	if isDAG {
		//worker.CurChain.Txpool = GenerateTxPool()
		//for len(worker.CurChain.Txpool.TxQueue) != params.TotalDataSize {
		//	time.Sleep(time.Second)
		//}
		go worker.MineDAGBlock()
		worker.TcpListen()
	} else {
		if nid == 0 {
			go worker.Propose()
			worker.TcpListen()
		} else {
			worker.TcpListen()
		}
	}

}

func data2tx(data []string, nonce uint64) (*core.Transaction, bool) {
	if data[6] == "0" && data[7] == "0" && len(data[3]) > 16 && len(data[4]) > 16 && data[3] != data[4] {
		val, ok := new(big.Int).SetString(data[8], 10)
		if !ok {
			log.Panic("new int failed\n")
		}
		tx := core.NewTransaction(data[3][2:], data[4][2:], val, nonce)
		return tx, true
	}
	return &core.Transaction{}, false
}

func GenerateTxPool() *core.TxPool {

	//params.FileInput, params.TotalDataSize, params.BatchSize

	txPool := core.TxPool{
		TxQueue:   make([]*core.Transaction, 0),
		RelayPool: make(map[uint64][]*core.Transaction),
	}

	dataTotalNum := 0

	dataPath := "D:/Projects/Blockchain/block-emulator/dataset/0to999999_BlockTransaction.csv"
	//_, fileName, _, _ := runtime.Caller(1)
	//dataPath := path.Join(path.Dir(fileName), filePath)
	fmt.Println(dataPath)

	txfile, err := os.Open(dataPath)
	if err != nil {
		log.Panic(err)
	}
	defer txfile.Close()
	reader := csv.NewReader(txfile)
	txlist := make([]*core.Transaction, 0) // save the txs in this epoch (round)

	for {
		data, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Panic(err)
		}
		if tx, ok := data2tx(data, uint64(dataTotalNum)); ok {
			txlist = append(txlist, tx)
			dataTotalNum++
		}

		if dataTotalNum == params.TotalDataSize {
			break
		}
	}

	for _, tx := range txlist {
		txPool.TxQueue = append(txPool.TxQueue, tx)
	}
	return &txPool

}
