package DAGBlock

import (
	"blockEmulator/core"
	"fmt"
	"os"
	"time"
)

func ChainAddBlock(Name string, References []string, p *Phantom) *core.Block {

	//create this block
	thisBlockDagData := core.DAGStructureData{
		Name:          Name,
		Score:         -1,
		SizeOfPastSet: -1,
		Prev:          make(map[string]*core.Block),
		Blue:          make(map[string]bool),
	}

	thisBlock := core.Block{
		Header: &core.BlockHeader{
			ParentBlockHash: nil,
			StateRoot:       nil,
			TxRoot:          nil,
			Number:          0,
			Time:            time.Time{},
			Miner:           0,
			DagData:         nil,
		},
		Body: nil,
		Hash: []byte(Name),
	}
	thisBlock.Header.DagData = &thisBlockDagData

	//为next添加新区块的键值对
	p.Next[Name] = make(MapBlocks)
	//add references
	for _, Reference := range References {
		prev, ok := p.G[Reference]
		if ok {
			thisBlock.Header.DagData.Prev[Reference] = prev
			if len(p.Next[prev.Header.DagData.Name]) == 0 {
				p.Next[prev.Header.DagData.Name] = make(MapBlocks)
			}
			p.Next[prev.Header.DagData.Name][Name] = &thisBlock
		} else {
			fmt.Println("chainAddBlock(): error! block reference invalid. block name =", Name, " references=", Reference)
			os.Exit(-1)
		}
	}

	thisBlock.Header.DagData.SizeOfPastSet = p.SizeOfPastSet(&thisBlock)

	//add this block to the chain
	p.G[Name] = &thisBlock
	return &thisBlock
}
