package main

import (
	"GoFrp/v1/server"
	"GoFrp/v1/slave"
	"GoFrp/v1/svcContext"
	"GoFrp/v1/util"
	"flag"
	"sync"
)

var mode = flag.String("m", "server", "run mode")
var port = flag.Int("p", 10000, "server bind port")
var serverHost = flag.String("h", "0.0.0.0", "server host")
var bindHost = flag.String("lh", "localhost", "local host")
var bindPort = flag.Int("lp", 443, "local bind port")
var password = flag.String("pwd", "12345678", "password for connect")

func main() {

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
		server.Start(svc)
	} else {
		slave.Start(svc)
	}
}
