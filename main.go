package main

import (
	"awesomeProject2/v2/server"
	"awesomeProject2/v2/slave"
	"flag"
)

var mode = flag.String("m", "server", "run mode")
var port = flag.Int("p", 10000, "data port")
var cmdPort = flag.Int("cp", 10001, "cmd port")
var remoteCmdHost = flag.String("rch", "localhost:10001", "remote cmd host")
var localHost = flag.String("lh", "192.168.0.65:8001", "remote cmd host")

func main() {
	flag.Parse()

	if *mode == "server" {
		go server.ListenServer(*port)
		server.ListenTunnelServer(*cmdPort)
	} else {
		slave.Start(*remoteCmdHost, *localHost)
	}
}
