package node

import (
	"github.com/maslow/xmemcache/config"
	"github.com/maslow/xmemcache/hasher"
	"strconv"
	"strings"
)

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
		for i := 0; i < config.GetCopyCount(); i++ {
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

//初始化，根据配置生成物理与虚拟节点
func (n *Nodes) Init(config []string) {
	n.generateRealNodeSet(config)
	n.generateVirtualNodeSet()
}

// 根据数据值导航到物理节点
func (n *Nodes) To(key string) (addr string) {
	vnhi := n.toVirtualNode(key)
	return n.toRealNode(vnhi)
}

// 根据数据键值得到所属的虚拟节点的hash索引
func (n *Nodes) toVirtualNode(key string) (vnhi uint32) {
	hashIndex := hasher.GetHashValue(key)

	vnhi = 0xFFFFFFFF
	min := vnhi
	for k, _ := range n.virtualNodes {
		if hashIndex <= k && vnhi > k {
			vnhi = k
		}
		if min > k {
			min = k
		}
	}
	if 0xFFFFFFFF == vnhi {
		if _, ok := n.virtualNodes[vnhi]; !ok {
			vnhi = min
		}
	}

	return vnhi
}

// 根据虚拟节点的hash索引得到物理节点地址
func (n *Nodes) toRealNode(vnhi uint32) (addr string) {
	vnode := n.virtualNodes[vnhi]
	str := strings.Split(vnode, "#")
	return str[0]
}
