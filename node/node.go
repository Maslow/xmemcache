package node

import (
	"errors"
	"fmt"
	"github.com/maslow/xmemcache/config"
	"github.com/maslow/xmemcache/hasher"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	NODE_STATUS_ACTIVE   = 0
	NODE_STATUS_DISABLED = 1
)

const SEPARATOR = "#"

type Nodes struct {
	mutex        sync.Mutex
	realNodes    map[string]uint8
	virtualNodes map[uint32]string
}

// Add virtual nodes by a given real node addr. see Add()
func (n *Nodes) addVirtualNodes(addr string) {
	for i := 0; i < config.GetCopyCount(); i++ {
		vnk := addr + SEPARATOR + strconv.Itoa(i)
		hashValue := hasher.GetHashValue(vnk)
		n.virtualNodes[hashValue] = vnk
	}
}

func (n *Nodes) removeVirtualNodes(addr string) {
	for i := 0; i < config.GetCopyCount(); i++ {
		vnk := addr + SEPARATOR + strconv.Itoa(i)
		hashValue := hasher.GetHashValue(vnk)
		delete(n.virtualNodes, hashValue)
	}
}

// Add a real node
func (n *Nodes) Add(addr string) error {
	if _, ok := n.realNodes[addr]; ok {
		return errors.New("the addr being added is already exist.")
	}
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.realNodes[addr] = NODE_STATUS_ACTIVE
	n.addVirtualNodes(addr)
	return nil
}

// Remove a real node
func (n *Nodes) Remove(addr string) {
	n.mutex.Lock()
	n.mutex.Unlock()
	n.removeVirtualNodes(addr)
	delete(n.realNodes, addr)
}

// Disable a real node
func (n *Nodes) Disable(addr string) error {
	if _, ok := n.realNodes[addr]; !ok {
		return errors.New("the addr being disabled is not exist.")
	}
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.removeVirtualNodes(addr)
	n.realNodes[addr] = NODE_STATUS_DISABLED
	return nil
}

// Enable a real node
func (n *Nodes) Enable(addr string) error {
	if _, ok := n.realNodes[addr]; !ok {
		return errors.New("the addr being enabled is not exist.")
	}
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.addVirtualNodes(addr)
	n.realNodes[addr] = NODE_STATUS_ACTIVE
	return nil
}

// Generate virtual nodes
func (n *Nodes) generateVirtualNodes() {
	n.virtualNodes = make(map[uint32]string)
	for k, _ := range n.realNodes {
		for i := 0; i < config.GetCopyCount(); i++ {
			vnk := k + SEPARATOR + strconv.Itoa(i)
			hashValue := hasher.GetHashValue(vnk)
			n.virtualNodes[hashValue] = vnk
		}
	}
}

// Generate real nodes
func (n *Nodes) generateRealNodes(servers []string) {
	n.realNodes = make(map[string]uint8)
	for _, v := range servers {
		n.realNodes[v] = NODE_STATUS_ACTIVE
	}
}

// Initialize , generate virtual nodes and real nodes
func (n *Nodes) Init(config []string) {
	n.generateRealNodes(config)
	n.generateVirtualNodes()
}

// Route to the node ip by the key
func (n *Nodes) To(key string) (addr string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	vnhi := n.toVirtualNode(key)
	return n.toRealNode(vnhi)
}

// Route to the related virtual node by the key
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

// Route to the real node by virtual node's hash-index
func (n *Nodes) toRealNode(vnhi uint32) (addr string) {
	vnode := n.virtualNodes[vnhi]
	str := strings.Split(vnode, SEPARATOR)
	return str[0]
}

// Manager the health of nodes
func (n *Nodes) Doctor() {
	limiter := time.Tick(time.Second * 5)
	for {
		<-limiter
		for addr := range n.realNodes {
			lnr, err := net.DialTimeout("tcp", addr, time.Second)

			if n.realNodes[addr] == NODE_STATUS_ACTIVE {
				i := 0
				for nil != err && i < 3 {
					lnr, err = net.DialTimeout("tcp", addr, time.Second)
					fmt.Printf("Retry to dialing the node named %s for %d times.\n", addr, i)
					i++
				}
				if i >= 3 && nil != err {
					n.Disable(addr)
					fmt.Printf("Disable the %s node.\n", addr)
				} else {
					lnr.Close()
				}
			} else {
				if nil == err {
					n.Enable(addr)
					lnr.Close()
					fmt.Printf("Enable the %s node.\n", addr)
				}
			}
		}
	}
}
