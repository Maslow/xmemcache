package main

import (
	"fmt"
	"github.com/maslow/xmemcache/config"
	"github.com/maslow/xmemcache/node"
	"github.com/maslow/xmemcache/protocol"
	"net"
	"os"
	"time"
)

var nodes *node.Nodes

func main() {
	nodes = new(node.Nodes)
	nodes.Init(config.GetServers())
	addr := config.GetListenAddr()
	lnr, e := net.Listen("tcp", addr)
	if nil != e {
		fmt.Fprint(os.Stderr, e)
	}
	defer lnr.Close()

	go nodes.Doctor()

	for {
		conn, e := lnr.Accept()
		if nil != e {
			fmt.Fprint(os.Stderr, e)
		}
		go deal(conn)
	}

}

func deal(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if nil != err {
			return
		}

		packet := new(protocol.Packet)
		if false == packet.Parse(buf[0:n]) {
			fmt.Fprint(os.Stderr, "Could not parse key from request")
			conn.Write([]byte("Error"))
			continue
		}
		key := packet.GetKey()
		ip := nodes.To(string(key))

		fmt.Printf("%s : ", key)
		fmt.Println(ip)

		mconn, err := net.DialTimeout("tcp", ip, time.Second)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			// TODO send [Not Found] response to client.
			continue
		}

		m, err := mconn.Write(buf[:n])
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			continue
		}
		m, err = mconn.Read(buf)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			continue
		}
		mconn.Close()
		conn.Write(buf[:m])
	}
}
