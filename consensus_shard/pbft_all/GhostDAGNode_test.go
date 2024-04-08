package pbft_all

import (
	"blockEmulator/message"
	"blockEmulator/params"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

func DeleteFolder(folderPath string) {
	// 检查文件夹是否存在
	if _, err := os.Stat(folderPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("The folder does not exist, so it cannot be deleted.")
		} else {
			fmt.Println("An error occurred while checking the folder:", err)
		}
		return
	}

	// 删除文件夹
	err := os.RemoveAll(folderPath)
	if err != nil {
		fmt.Println("An error occurred while deleting the folder:", err)
	} else {
		fmt.Println("The folder has been successfully deleted.")
	}
}

func initConfig_test(nid, nnm, sid, snm uint64) *params.ChainConfig {
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

func TestGhostDAGNodeMining(t *testing.T) {
	DeleteFolder("./record")

	worker := NewPbftNode(0, 0, initConfig_test(0, 4, 0, 2), params.CommitteeMethod[3], true)

	go worker.MineDAGBlock()
	worker.TcpListen()

	DeleteFolder("./record")

}

func TestGhostDAGNodeReceiving(t *testing.T) {
	DeleteFolder("./record")

	worker := NewPbftNode(0, 0, initConfig_test(0, 4, 0, 2), params.CommitteeMethod[3], true)

	var msgs [][]byte

	for i := 0; i < 10; i++ {
		// implement interface to generate propose
		block := worker.CurChain.GenerateBlock()
		r := &message.Request{
			RequestType: message.BlockRequest,
			ReqTime:     time.Now(),
		}
		r.Msg.Content = block.Encode()
		//生成区块，打包进消息
		digest := getDigest(r)
		worker.requestPool[string(digest)] = r
		worker.pl.Plog.Printf("S%dN%d put the request into the pool ...\n", worker.ShardID, worker.NodeID)

		DAGSendBlockMsg := message.DAGSendBlock{
			RequestMsg: r,
			Digest:     digest,
		}
		// marshal and broadcast
		DAGSendBlockMsgBytes, err := json.Marshal(DAGSendBlockMsg)
		if err != nil {
			log.Panic()
		}
		msg_send := message.MergeMessage(message.CGhostDAGSendBlock, DAGSendBlockMsgBytes)
		msgs = append(msgs, msg_send)
	}

	for _, v := range msgs {
		worker.handleMessage(v)
	}

	//worker.TcpListen()

	DeleteFolder("./record")

}
