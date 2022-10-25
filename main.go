package main

import (
	"GoFrp/config"
	"GoFrp/multi_wire/client"
	"GoFrp/multi_wire/server"
	"GoFrp/multi_wire/svcContext"
	"GoFrp/multi_wire/util"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
)

var filepath = flag.String("c", "./config.json", "config file path")

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.Parse()

	configs, err := config.ReadConfig(*filepath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for _, config := range configs {
		svc := &svcContext.SVCContext{
			ApplyNewDataTunChan: nil,
			TaskMap:             sync.Map{},
			ServerPort:          config.ServerPort,
			ServerHost:          config.ServerHost,
			BindHost:            config.BindHost,
			BindPort:            config.BindPort,
			Password:            util.ParsePassword(config.Password),
		}
		if config.Mode == "server" {
			go server.Start(svc)
		} else {
			go client.Start(svc)
		}
	}

	go testUDP()

	<-ctx.Done()
	fmt.Println("bye")
}

func testUDP() {
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 10000,
	})
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	fmt.Println("listen success")
	defer listen.Close()

	var sourceaddr net.UDPAddr

	for {
		var data [10240]byte
		n, addr, err := listen.ReadFromUDP(data[:]) // 接收数据
		if err != nil {
			fmt.Println("read udp failed, err:", err)
			continue
		}
		fmt.Println("read udp data")
		if addr.IP.IsLoopback() {

			_, err = listen.WriteToUDP(data[:n], &sourceaddr) // 发送数据
			if err != nil {
				fmt.Println("write to udp failed, err:", err)
				continue
			}

		} else {

			sourceaddr = *addr

			naddr := &net.UDPAddr{
				IP:   net.IPv4(127, 0, 0, 1),
				Port: 5900,
			}
			_, err = listen.WriteToUDP(data[:n], naddr) // 发送数据
			if err != nil {
				fmt.Println("write to udp failed, err:", err)
				continue
			}
		}

		fmt.Printf("data:%v addr:%v count:%v\n", string(data[:n]), addr, n)
		_, err = listen.WriteToUDP(data[:n], addr) // 发送数据
		if err != nil {
			fmt.Println("write to udp failed, err:", err)
			continue
		}
	}
}
