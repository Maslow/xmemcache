package main

import (
	"fmt"
	"github.com/maslow/xmemcache/config"
	"github.com/maslow/xmemcache/node"
	"github.com/maslow/xmemcache/protocol"
	"net"
	"os"
)

var nodes *node.Nodes

func main() {
	nodes = new(node.Nodes)
	nodes.Init(config.GetServers())

	lnr, err := net.Listen("tcp", config.GetListenAddr())
	if nil != err {
		fmt.Fprint(os.Stderr, err)
	}
	fmt.Fprintf(os.Stdout, "Start listening on %s\n", config.GetListenAddr())
	defer lnr.Close()
	for {
		conn, err := lnr.Accept()
		if nil != err {
			fmt.Fprint(os.Stderr, err)
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
			continue
		}
		key := packet.GetKey()
		ip := nodes.To(string(key))
		fmt.Printf("%s : ", key)
		fmt.Println(ip)

		//TODO  添加超时处理
		//TODO  实现连接池，避免重复连接
		mconn, err := net.Dial("tcp", ip)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
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
