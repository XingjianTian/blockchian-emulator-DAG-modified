// definition of node and shard

package shard

import (
	"fmt"
	"strconv"
)

type Node struct {
	NodeID  uint64
	ShardID uint64
	IPaddr  string
}

func (n *Node) PrintNode() string {
	v := []interface{}{
		n.NodeID,
		n.ShardID,
		n.IPaddr,
	}
	fmt.Printf("%v\n", v)
	return "S" + strconv.FormatUint(n.ShardID, 10) + "N" + strconv.FormatUint(n.NodeID, 10)
}
