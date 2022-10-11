package main

import (
	"GoFrp/v1/server"
	"GoFrp/v1/server/svcContext"
	"GoFrp/v1/slave"
	"flag"
	"sync"
)

var mode = flag.String("m", "server", "run mode")
var port = flag.Int("p", 10000, "data port")
var cmdPort = flag.Int("cp", 10001, "cmd port")
var remoteCmdHost = flag.String("rch", "localhost:10001", "remote cmd host")
var localHost = flag.String("lh", "192.168.0.65:8001", "remote cmd host")

func main() {
	flag.Parse()

	if *mode == "server" {
		svc := &svcContext.SVCContext{
			ApplyNewDataTunChan: nil,
			TaskMap:             sync.Map{},
			ServerPort:          *port,
		}

		server.StartServer(svc)
	} else {
		slave.Listen(*remoteCmdHost, *localHost)
	}
}
