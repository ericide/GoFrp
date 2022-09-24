package server

import (
	"fmt"
	"log"
	"net"
)

func Start(bindPort int, tunnelPort int) {
	tunnelServer := &TunnelServer{
		Port: tunnelPort,
	}
	go tunnelServer.Start()

	startExtServer(bindPort, tunnelServer)
}

func startExtServer(port int, tunnelServer *TunnelServer) {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		log.Println("Error listening", err.Error())
		return //终止程序
	}
	// 监听并接受来自客户端的连接
	log.Println("start Listening ext conn", port)
	for {
		conn, err := listener.Accept()
		log.Printf("Ext accepted new ext connection: %v \n", conn)

		if err != nil {
			log.Println("Error accepting", err.Error())
			return // 终止程序
		}
		tunnelServer.AddNewTunnelTask(&conn)
	}
}
