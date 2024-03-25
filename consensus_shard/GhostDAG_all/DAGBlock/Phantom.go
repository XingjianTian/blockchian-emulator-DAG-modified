package DAGBlock

import (
	"blockEmulator/core"
	"github.com/deckarep/golang-set"
	"github.com/heimdalr/dag"
	"time"
)

type Block = core.Block
type BlockHeader = core.BlockHeader

// Phantom 结构体实现了 DAG 接口
type Phantom struct {
	G           *dag.DAG   // 有向图对象
	Leaves      mapset.Set // 所有叶子节点的集合
	Coloring    mapset.Set // 所有蓝色区块的集合
	GenesisHash []byte     // 创世区块的哈希
	k           int        // 蓝色反锥体的最大大小
}

// NewPhantom 创建并初始化一个新的 Phantom 实例
func NewPhantom(k int) *Phantom {
	return &Phantom{
		G:           dag.NewDAG(),
		Leaves:      mapset.NewSet(),
		Coloring:    mapset.NewSet(),
		GenesisHash: nil,
		k:           k,
	}
}

// Init 初始化 Phantom DAG
func (p *Phantom) Init() error {
	// 在这里实现初始化逻辑，例如添加创世区块
	// 示例：创建创世区块并将其添加到 DAG
	genesisBlock := &Block{
		Header: &BlockHeader{
			Number: 0,
			Time:   time.Unix(0, 0),
		},
		Body: nil,
		Hash: nil, // 创世区块的哈希通常设置为特定的值，如空哈希或零哈希
	}
	p.GenesisHash = genesisBlock.Hash // 假设 Block 结构体有一个方法来计算哈希
	p.G.AddVertexByHash(genesisBlock.Hash, genesisBlock)
	p.Leaves.Add(genesisBlock.Hash)
	return nil
}

// Contains 检查 DAG 中是否包含指定哈希的区块
func (p *Phantom) Contains(hash []byte) bool {
	_, err := p.G.GetVertexByHash(hash)
	return err == nil
}

// Get 获取 DAG 中指定哈希的区块
func (p *Phantom) Get(hash []byte) *Block {
	node, _ := p.G.GetVertexByHash(hash)
	if node == nil {
		return nil
	}
	blockData, ok := node.(*Block) //断言
	if !ok {
		return nil
	}
	return blockData
}

// Iterator 返回 DAG 中所有区块的迭代器
func (p *Phantom) Iterator() <-chan *Block {
	blocks := make(chan *Block)
	go func() {
		for _, node := range p.G.GetVertices() {
			blockData, ok := node.(*Block)
			if !ok {
				continue
			}
			blocks <- blockData
		}
		close(blocks)
	}()
	return blocks
}

// Len 获取 DAG 中区块的数量
func (p *Phantom) Len() int {
	return p.G.GetOrder()
}

// String 获取 DAG 的字符串表示
func (p *Phantom) String() string {
	return p.G.String()
}

// GetVirtualBlockParents 获取虚拟块的父块哈希集合
func (p *Phantom) GetVirtualBlockParents() mapset.Set {
	// 在这里实现获取虚拟块父块的逻辑
	leaves := p.G.GetLeaves()
	for hash, _ := range leaves {
		p.Leaves.Add([]byte(hash))
	}
	return p.Leaves
}

func (p *Phantom) GetPast(hash []byte) (map[string]interface{}, error) {
	return p.G.GetAncestorsByHash(hash)
}

func (p *Phantom) GetFuture(hash []byte) (map[string]interface{}, error) {
	return p.G.GetDescendantsByHash(hash)
}

// GetAnticone 返回给定区块的反锥体。
func (p *Phantom) GetAnticone(hash []byte) mapset.Set {
	// 获取区块的过去集合和未来集合
	past, _ := p.GetPast(hash)
	future, _ := p.GetFuture(hash)

	// 合并过去集合和未来集合
	blockConeUnion := mapset.NewSet(hash)
	for ph, _ := range past {
		blockConeUnion.Add(ph)
	}
	for fh, _ := range future {
		blockConeUnion.Add(fh)
	}
	// 获取所有区块的集合
	allBlocksHashes := p.G.GetVertices() // 假设GetAllBlocks能够返回所有区块的集合
	allBlocksUnion := mapset.NewSet()
	for h, _ := range allBlocksHashes {
		allBlocksUnion.Add(h)
	}
	// 计算反锥体
	antiConeSet := allBlocksUnion.Difference(blockConeUnion)

	return antiConeSet
}

// Add 向DAG中添加一个区块
func (p *Phantom) Add(block *Block) {
	hash := block.Encode() // 假设Block有一个Encode方法来序列化区块

	// 获取区块的父块哈希值
	parents := block.Header.DAGParentBlockHashes

	// 将区块添加到DAG中
	p.G.AddVertexByHash(hash, block)
	for _, parentHash := range parents {
		p.G.AddEdgeByHashes(parentHash, hash)
		p.Leaves.Remove(parentHash)
	}

	// 更新叶子节点集合
	p.Leaves.Add(hash)

	// 更新区块的着色以及拓扑排序
	//p.UpdateColoringIncrementally(hash)
	//p._update_topological_order_incrementally(hash)
}

/*
func (p *Phantom) UpdateColoringIncrementally(hash []byte) {
	// 计算区块的反锥体
	anticones := p.GetAnticone(hash)

	// 计算所有可能的着色集合的幂集
	var colorings []mapset.Set
	for _, block := range p.G.GetVertices() {
		colorings = append(colorings, mapset.NewSet(block))
	}
	for i := 1; i < len(colorings); i++ {
		for j := 0; j < i; j++ {
			colorings = append(colorings, colorings[i-1].Union(colorings[j]))
		}
	}

	// 寻找最大的有效着色集合
	var maxColoring mapset.Set
	var maxColoringBAC map[string]mapset.Set
	for _, coloring := range colorings {
		coloringBAC := make(map[string]mapset.Set)
		for block, anticone := range anticones.ToSlice() {
			blueAnticone := anticone.Intersect(coloring)
			if _, exists := coloringBAC[block]; !exists {
				coloringBAC[block] = blueAnticone
			} else {
				coloringBAC[block] = coloringBAC[block].Union(blueAnticone)
			}
			if len(blueAnticone) > p.k {
				break
			}
		}
		if len(coloring) > len(maxColoring) {
			maxColoring = coloring
			maxColoringBAC = coloringBAC
		}
	}

	// 更新着色集合
	p.Coloring = maxColoring

	// 根据新的着色更新每个区块的蓝色反锥体
	p.muDAG.Lock()
	for block, blueAnticone := range maxColoringBAC {
		if node, exists := p.G.GetNodeByHash(block); exists {
			node.Item.(*Block).BlueAnticone = blueAnticone.List() // 假设Block结构体有一个BlueAnticone字段
		}
	}
	p.muDAG.Unlock()
}
*/
