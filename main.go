package main

import (
	"GoFrp/multi_wire/client"
	"GoFrp/multi_wire/server"
	"GoFrp/multi_wire/svcContext"
	"GoFrp/multi_wire/util"
	"context"
	"flag"
	"fmt"
	"sync"
)

var mode = flag.String("m", "server", "run mode")
var port = flag.Int("p", 10000, "server bind port")
var serverHost = flag.String("h", "0.0.0.0", "server host")
var bindHost = flag.String("lh", "localhost", "local host")
var bindPort = flag.Int("lp", 443, "local bind port")
var password = flag.String("pwd", "12345678", "password for connect")

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.Parse()

	svc := &svcContext.SVCContext{
		ApplyNewDataTunChan: nil,
		TaskMap:             sync.Map{},
		ServerPort:          *port,
		ServerHost:          *serverHost,
		BindHost:            *bindHost,
		BindPort:            *bindPort,
		Password:            util.ParsePassword(*password),
	}

	if *mode == "server" {
		go server.Start(svc)
	} else {
		go client.Start(svc)
	}

	<-ctx.Done()
	fmt.Println("bye")
}
