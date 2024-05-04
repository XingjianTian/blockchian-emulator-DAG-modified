package committee

import (
	"blockEmulator/message"
)

type CommitteeModule interface {
	HandleBlockInfo(*message.BlockInfoMsg)
	MsgSendingControl(isDAG bool)
	HandleOtherMessage([]byte)
}
