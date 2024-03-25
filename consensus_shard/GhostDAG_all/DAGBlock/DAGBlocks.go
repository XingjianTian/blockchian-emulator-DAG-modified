package DAGBlock

import mapset "github.com/deckarep/golang-set"

type DAGBlocks interface {

	// 初始化DAG
	Init() error

	// 检查DAG中是否包含指定全局ID的区块
	Contains(hash []byte) bool

	// 获取DAG中指定全局ID的区块
	Get(hash []byte) *Block

	// 返回DAG中所有区块的迭代器
	Iterator() <-chan *Block

	// 获取DAG中区块的数量
	Len() int

	// 获取虚拟块的父块全局ID集合
	GetVirtualBlockParents() mapset.Set

	// 向DAG中添加一个区块
	Add(block *Block) error

	// 判断区块a是否在区块b之前
	IsABeforeB(a, b []byte) bool

	// 获取指定全局ID区块在主子DAG中的深度
	GetDepth(hash []byte) int
}
