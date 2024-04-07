package pbft_all

import "blockEmulator/message"

type ExtraOpInConsensus interface {
	// mining / message generation
	HandleinPropose() (bool, *message.Request)
	// checking
	HandleinPrePrepare(*message.PrePrepare) bool
	// nothing necessary
	HandleinPrepare(*message.Prepare) bool
	// confirming
	HandleinCommit(*message.Commit) bool
	// do for need
	HandleReqestforOldSeq(*message.RequestOldMessage) bool
	// do for need
	HandleforSequentialRequest(*message.SendOldMessage) bool

	//Ghost DAG
	HandleinDAGMineBlock() (bool, *message.Request)
	//生成区块，打包进消息
	// confirming
	HandleinDAGReceiveBlock(receivedBlock *message.DAGSendBlock) bool
}

type OpInterShards interface {
	// operation inter-shards
	HandleMessageOutsidePBFT(message.MessageType, []byte) bool
}
