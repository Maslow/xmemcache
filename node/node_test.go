package node

import (
	"fmt"
	"github.com/maslow/xmemcache/config"
	"testing"
)

func Test_Add(t *testing.T) {
	t.Error("TODO please write the test code")
}

func Test_Remove(t *testing.T) {
	t.Error("TODO please write the test code")
}

func Test_To(t *testing.T) {
	servers := config.GetServers()
	nodes := new(Nodes)
	nodes.Init(servers)
	key := "rocky52"
	if nodes.To(key) != nodes.To(key) {
		t.Error("Error: these two result should be equal to each other")
	}
	rn := nodes.To(key)
	i := 0
	for ; i < len(servers); i++ {
		if rn == servers[i] {
			break
		}
	}
	if i >= len(servers) {
		t.Errorf("the result is not in the real node list provided by config: %s", rn)
	}
}

//测试虚拟节点分散程度
func Test_hitRate(t *testing.T) {
	servers := config.GetServers()
	nodes := new(Nodes)
	nodes.Init(servers)

	rnCount := uint32(len(nodes.realNodes))
	var max uint32 = 0xFFFFFFFF

	T := make([]uint32, rnCount)
	for i := uint32(0); i < rnCount; i++ {
		T[i] = (i + 1) * (max / rnCount)
	}

	C := make([]uint32, rnCount)

	for k, _ := range nodes.virtualNodes {
		for ii := uint32(0); ii < rnCount; ii++ {
			if k <= T[ii] {
				C[ii]++
			}
		}
	}

	for j := rnCount - 1; j > uint32(0); j-- {
		C[j] = C[j] - C[j-1]
	}
	fmt.Println("虚拟节点分散情况:", C)
}
