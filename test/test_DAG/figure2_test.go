package test_DAG

import (
	"blockEmulator/consensus_shard/GhostDAG_all/DAGBlock"
	"bytes"
	"fmt"
	"testing"
)

func chainFig2Initialize(p *DAGBlock.Phantom) DAGBlock.MapBlocks {

	//add blocks
	DAGBlock.ChainAddBlock("Genesis", []string{}, p)

	DAGBlock.ChainAddBlock("B", []string{"Genesis"}, p)
	DAGBlock.ChainAddBlock("C", []string{"Genesis"}, p)
	DAGBlock.ChainAddBlock("D", []string{"Genesis"}, p)
	DAGBlock.ChainAddBlock("E", []string{"Genesis"}, p)

	DAGBlock.ChainAddBlock("F", []string{"B", "C"}, p)
	DAGBlock.ChainAddBlock("G", []string{"C", "D"}, p)
	DAGBlock.ChainAddBlock("H", []string{"E"}, p)

	DAGBlock.ChainAddBlock("I", []string{"F", "D"}, p)
	DAGBlock.ChainAddBlock("J", []string{"B", "G", "E"}, p)
	DAGBlock.ChainAddBlock("K", []string{"D", "H"}, p)

	tips := p.FindTips(p.G)
	tipsName := p.LTPQ(tips, true) // LTPQ is not relevant here, I just use it to get Tips name.
	DAGBlock.ChainAddBlock("Virtual", tipsName, p)

	//p.Tips = tips
	return p.G
}

// Tests Algorithm 2 Ordering of the DAG, with the example on paper page 3 Fig.2
func TestFig2(t *testing.T) {

	var actual bytes.Buffer
	var expected string

	fmt.Println("\n-  BlockDAG Algorithms Simulation - Algorithm 2: Ordering of the DAG.   -")
	fmt.Println("-        The example on Fig.2         -")

	p := DAGBlock.NewPhantom(3)

	chain := chainFig2Initialize(p)

	fmt.Println("chainInitialize(): done. blocks=", len(chain)-1)

	orderedList := p.Order(chain, 3)

	// print the result of blue sets

	ltpq := p.LTPQ(chain, true)

	fmt.Print("blue set selection done. blue blocks = ")
	nBlueBlocks := 0
	actual.Reset()
	for _, name := range ltpq {
		block := chain[name]
		if p.IsBlueBlock(block) == true {
			if name == "Genesis" || name == "Virtual" {
				actual.WriteString(fmt.Sprintf("(%s).", name[:1]))
			} else {
				actual.WriteString(name + ".")
			}

			nBlueBlocks++
		}
	}
	fmt.Println(actual.String(), "	total blue:", nBlueBlocks)

	expected = "(G).B.C.D.F.G.I.J.(V)."
	if actual.String() != expected {
		t.Errorf("blue selection test not matched. expected=%s", expected)
	}

	// print the result of ordered blocks of this chain

	fmt.Print("ordered chain blocks = ")
	actual.Reset()
	for _, name := range orderedList {
		if name == "Genesis" || name == "Virtual" {
			actual.WriteString(fmt.Sprintf("(%s).", name[:1]))
		} else {
			actual.WriteString(name + ".")
		}
	}
	fmt.Println(actual.String())

	expected = "(G).B.C.D.F.E.G.J.I.H.K.(V)."
	if actual.String() != expected {
		t.Errorf("block ordering test not matched. expected=%s", expected)
	}
}

func BenchmarkBlockOrdering(b *testing.B) {

	//for i := 0; i < b.N; i++ {
	//chain := chainFig2Initialize()
	//Order(chain, 3)
	//}
}
