package node

import (
	"github.com/maslow/xmemcache/hasher"
	//"github.com/maslow/xmemcache/config"
	"strconv"
)

//默认<复制个数>
const defaultCopyCount = 50

const (
	NODE_STATUS_ACTIVE   = 0
	NODE_STATUS_INACTIVE = 1
)

type Nodes struct {
	realNodes    map[string]uint8
	virtualNodes map[uint32]string
}

//生成虚拟节点
func (n *Nodes) generateVirtualNodeSet() {
	n.virtualNodes = make(map[uint32]string)
	for k, _ := range n.realNodes {
		for i := 0; i < defaultCopyCount; i++ {
			vnKey := k + "#" + strconv.Itoa(i)
			hashValue := hasher.GetHashValue(vnKey)
			n.virtualNodes[hashValue] = vnKey
		}
	}
}

//生成物理节点
func (n *Nodes) generateRealNodeSet(servers []string) {
	n.realNodes = make(map[string]uint8)
	for _, v := range servers {
		n.realNodes[v] = NODE_STATUS_ACTIVE
	}
}

//初始化
func (n *Nodes) Init(config []string) {
	n.generateRealNodeSet(config)
	n.generateVirtualNodeSet()
}

// TODO
func (n *Nodes) To(key string) (addr string) {
	return ""
}

// TODO
func (n *Nodes) toVirtualNode(key string) (hashValue uint32) {
	return 0
}

// TODO
func (n *Nodes) toRealNode() (addr string) {
	return ""
}
