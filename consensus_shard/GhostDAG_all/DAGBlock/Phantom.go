package DAGBlock

import (
	"blockEmulator/core"
	"fmt"
	"os"
	"sort"
	"time"
)

type MapBlocks = map[string]*core.Block

// Phantom 结构体做一层封装
type Phantom struct {
	G     MapBlocks // 有向图对象集合
	Tips  MapBlocks // 所有叶子节点的集合
	Blues MapBlocks // 所有蓝色区块的集合

	Next        map[string]MapBlocks // Block don't have this info, it comes from the analysis of existing chain
	GenesisHash core.BlockHash       // 创世区块的哈希
	k           int                  // 蓝色反锥体的最大大小
}

// NewPhantom 创建并初始化一个新的 Phantom 实例
func NewPhantom(k int) *Phantom {

	return &Phantom{
		G:           make(MapBlocks),
		Tips:        make(MapBlocks),
		Blues:       make(MapBlocks),
		Next:        make(map[string]MapBlocks),
		GenesisHash: nil,
		k:           k,
	}
}

// Init 初始化 Phantom DAG
func (p *Phantom) Init() error {
	// 在这里实现初始化逻辑，例如添加创世区块
	// 示例：创建创世区块并将其添加到 DAG
	genesisBlock := &core.Block{
		Header: &core.BlockHeader{
			Number: 0,
			Time:   time.Unix(0, 0),
		},
		Body: nil,
		Hash: nil, // 创世区块的哈希通常设置为特定的值，如空哈希或零哈希
	}
	p.GenesisHash = genesisBlock.Hash // 假设 Block 结构体有一个方法来计算哈希
	p.AddBlockByHash(genesisBlock.Hash, genesisBlock)
	//p.Tips.Add(genesisBlock.Hash)
	return nil
}

func (p *Phantom) AddBlockByStringName(blockHash string, block *core.Block) {

	//1.添加block，设置prev(past)
	//2. 给prev block添加next(future)

	block.Header.DagData.SizeOfPastSet = p.SizeOfPastSet(block)

	//add this block to the chain
	p.G[blockHash] = block
}
func (p *Phantom) AddBlockByHash(blockHash []byte, block *core.Block) {
	blockHashStr := string(blockHash)
	p.AddBlockByStringName(blockHashStr, block)
}

func (p *Phantom) FindTips(Graph MapBlocks) MapBlocks {

	tips := make(MapBlocks)

checkNextBlock:
	for k, v := range Graph {

		if len(p.Next[k]) == 0 {
			tips[k] = v
		} else {
			for next := range p.Next[v.Header.DagData.Name] {
				if _, ok := Graph[next]; ok {
					continue checkNextBlock
				}
			}

			tips[v.Header.DagData.Name] = v
		}
	}

	return tips
}

// PastSet 方法用于计算给定区块的过去集（past set）。
func (p *Phantom) PastSet(B *core.Block, past MapBlocks) {

	for k, v := range B.Header.DagData.Prev {

		if _, ok := past[k]; !ok {
			p.PastSet(v, past)
		}
		past[k] = v
	}
}

// FutureSet 方法用于计算给定区块的未来集（future set）。
func (p *Phantom) FutureSet(B *core.Block, future MapBlocks) {
	for k, v := range p.Next[B.Header.DagData.Name] {

		if _, ok := future[k]; !ok {
			p.FutureSet(v, future)
		}
		future[k] = v
	}
}

func (p *Phantom) countBlue(Graph MapBlocks, tip *core.Block) int {

	var blueBlocks = 0
	for _, v := range Graph {
		if blue, ok := v.Header.DagData.Blue[tip.Header.DagData.Name]; ok {
			if blue {
				blueBlocks++
			}
		} else if v.Header.DagData.Name == "Genesis" {
			blueBlocks++
		}

	}

	return blueBlocks
}

func (p *Phantom) AntiCone(Graph MapBlocks, B *core.Block) MapBlocks {

	anticone := make(MapBlocks)

	past := make(MapBlocks)
	p.PastSet(B, past)

	future := make(MapBlocks)
	p.FutureSet(B, future)

	for name, block := range Graph {
		if _, ok := past[name]; ok {
			continue // block not in B's past
		}

		if _, ok := future[name]; ok {
			continue // block not in B's future
		}

		if name == B.Header.DagData.Name {
			continue // block not B
		}

		anticone[name] = block // then this block belongs to anticone
	}

	return anticone
}

func (p *Phantom) IsBlueBlock(B *core.Block) bool {

	if B == nil {
		return false
	}

	if B.Header.DagData.Name == "Genesis" {
		return true
	}

	for _, blue := range B.Header.DagData.Blue {
		if blue == true {
			return true
		}
	}

	return false
}

func (p *Phantom) SizeOfPastSet(B *core.Block) int {

	past := make(MapBlocks)
	p.PastSet(B, past)
	return len(past)
}

func (p *Phantom) LTPQ(Graph MapBlocks, ascending bool) []string {

	ltpq := make([]struct {
		Name          string
		SizeOfPastSet int
	}, len(Graph))

	i := 0
	for _, block := range Graph {
		ltpq[i].Name, ltpq[i].SizeOfPastSet = block.Header.DagData.Name, block.Header.DagData.SizeOfPastSet
		i++
	}

	/*
	 * The priority of a block is represented by the size of its past set; in case of ties, the block with lowest hash ID is chosen.
	 */
	sort.Slice(ltpq, func(i, j int) bool {
		return ltpq[i].SizeOfPastSet < ltpq[j].SizeOfPastSet || (ltpq[i].SizeOfPastSet == ltpq[j].SizeOfPastSet && ltpq[i].Name < ltpq[j].Name)
	})

	priorityQue := make([]string, len(ltpq))

	if ascending == true { // asc: little first
		for i := 0; i < len(ltpq); i++ {
			priorityQue[i] = ltpq[i].Name
		}
	} else { // desc: big first
		for i, j := len(ltpq)-1, 0; i >= 0; i, j = i-1, j+1 {
			priorityQue[j] = ltpq[i].Name
		}
	}

	return priorityQue
}

func (p *Phantom) CalcBlue(Graph MapBlocks, k int, topTip *core.Block) {

	if len(Graph) == 1 {
		if _, ok := Graph["Genesis"]; ok {
			return
		} else {
			fmt.Println("CalcBlue(): error! len(G)=1 but not Genesis block")
			os.Exit(-1)
		}
	} else if len(Graph) == 0 {
		fmt.Println("CalcBlue(): error! impossible to reach here. len(G)=0")
		os.Exit(-1)
	}

	// step 4
	tips := p.FindTips(Graph)
	if len(tips) == 0 {
		fmt.Println("calcBlue(): error! impossible! Tips Empty.")
		os.Exit(-1)
	}

	maxBlue := -1
	var Bmax *core.Block = nil

	// step 4'	 	Starting from the highest scoring tip (to be continued...  for the moment, I use reversed LTPQ.)
	var ltpq = p.LTPQ(tips, false)

	for _, name := range ltpq {
		tip := tips[name] // step 4'

		past := make(MapBlocks)
		p.PastSet(tip, past)

		//fmt.Println("calcBlue(): info of next recursive call - tip=", tip.Name, " past=", len(past))

		// step 5
		p.CalcBlue(past, k, tip)

		// step 6
		blueBlocks := p.countBlue(past, tip)
		if blueBlocks > maxBlue {
			maxBlue = blueBlocks
			Bmax = tip
		} else if blueBlocks == maxBlue {
			// tie-breaking
			/*
			 * Important Note: in some cases, the tie-breaking method will decide the final blue selection result! not always converge to maximum k-cluster SubDAG.
			 *				   research to be continued.
			 */
			if tip.Header.DagData.Name < Bmax.Header.DagData.Name {
				Bmax = tip
			}
		}
	}

	if Bmax == nil {
		fmt.Println("calcBlue(): error! impossible! Bmax=nil.")
		os.Exit(-1)
	}

	// step 7
	for _, v := range Graph {
		for name, blue := range v.Header.DagData.Blue {
			for _, tip := range tips {
				if blue == true && name != Bmax.Header.DagData.Name && name == tip.Header.DagData.Name {
					v.Header.DagData.Blue[name] = false // clear all other tips blue blocks, only keep the Bmax blue ones
				}
			}
		}
	}
	Bmax.Header.DagData.Blue[Bmax.Header.DagData.Name] = true // BLUEk(G) = BLUEk(Bmax) U {Bmax}

	// step 8
	anticoneBmax := p.AntiCone(Graph, Bmax)
	ltpq = p.LTPQ(anticoneBmax, true) // in 'some' topological ordering: LTPQ
	for _, name := range ltpq {

		B := anticoneBmax[name]

		// step 9
		nBlueAnticone := 0
		anticone := p.AntiCone(Graph, B)
		for _, block := range anticone {
			if p.IsBlueBlock(block) == true {
				nBlueAnticone++
			}
		}

		if nBlueAnticone <= k {
			B.Header.DagData.Blue[Bmax.Header.DagData.Name] = true // step 10
		}
	}

	// additional step: replace Blue[Bmax] with [topTip]
	for _, B := range Graph {
		if blue, ok := B.Header.DagData.Blue[Bmax.Header.DagData.Name]; ok && blue == true {
			B.Header.DagData.Blue[Bmax.Header.DagData.Name] = false
			B.Header.DagData.Blue[topTip.Header.DagData.Name] = true
		}
	}

}

func (p *Phantom) Intersection(A MapBlocks, B MapBlocks) MapBlocks {

	intersection := make(MapBlocks)

	for _, blockA := range A {
		for _, blockB := range B {
			if blockA.Header.DagData.Name == blockB.Header.DagData.Name {
				intersection[blockA.Header.DagData.Name] = blockA
			}
		}
	}

	return intersection
}

func (p *Phantom) Order(Graph MapBlocks, k int) []string {

	// step 2
	topo_queue := make([]*core.Block, 0)
	map_queue := make(MapBlocks) // store the blocks already in queue

	// step 3
	ordered_list := make([]string, 0)
	map_list := make(map[string]int) // store the order of the blocks
	var order = 0

	// step 4
	p.CalcBlue(Graph, k, Graph["Virtual"])

	// step 5
	topo_queue = append(topo_queue, Graph["Genesis"]) // FIFO Queue Push
	map_queue["Genesis"] = Graph["Genesis"]

	// step 6
	var B *core.Block = nil
	for len(topo_queue) > 0 {

		B = topo_queue[0]           // step 7.	FIFO Queue Pop
		topo_queue = topo_queue[1:] // discard top element of FIFO Queue

		// step 8
		ordered_list = append(ordered_list, B.Header.DagData.Name)
		map_list[B.Header.DagData.Name] = order
		order++

		// step 9
		childrenB := p.Next[B.Header.DagData.Name]

		// step 9'		in 'some' topological ordering: LTPQ
		ltpq := p.LTPQ(childrenB, true)

		for _, name := range ltpq {
			C := childrenB[name] // step 9"

			if p.IsBlueBlock(C) == true {
				// step 10
				pastC := make(MapBlocks)
				p.PastSet(C, pastC)

				anticoneB := p.AntiCone(Graph, B)
				intersection := p.Intersection(pastC, anticoneB)

				// step 10'		in 'some' topological ordering: LTPQ
				ltpq2 := p.LTPQ(intersection, true)

				for _, name := range ltpq2 {
					D := intersection[name] // step 10"

					if _, okk := map_list[D.Header.DagData.Name]; !okk {
						// step 11
						if _, ok3 := map_queue[D.Header.DagData.Name]; !ok3 { // the queue should avoid duplicating elements
							topo_queue = append(topo_queue, D) // FIFO Queue Push
							map_queue[D.Header.DagData.Name] = D
						}
					}
				}

				// step 12
				if _, ok4 := map_queue[C.Header.DagData.Name]; !ok4 { // the queue should avoid duplicating elements
					topo_queue = append(topo_queue, C) // FIFO Queue Push
					map_queue[C.Header.DagData.Name] = C
				}
			}
		}
	}

	return ordered_list
}
