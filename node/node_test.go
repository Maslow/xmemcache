package node

import (
	"fmt"
	"testing"
)

func Test_generateReadNodeSet(t *testing.T) {

}

func Test_generateVirtualNodeSet(t *testing.T) {

}

//测试虚拟节点分散程度
func Test_hitRate(t *testing.T) {
	config := []string{"12.108.111.20", "19.168.3.61", "192.168.33.52", "192.168.333.13", "192.68.33.14", "12.168.33.15", "192.16.3.16"}
	nodes := new(Nodes)
	nodes.Init(config)

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
