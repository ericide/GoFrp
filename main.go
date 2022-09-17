package main

import (
	"awesomeProject2/client"
	"awesomeProject2/server"
	"flag"
)

var mode = flag.String("m", "server", "run mode")
var port = flag.Int("p", 10000, "data port")
var cmdPort = flag.Int("cp", 10001, "cmd port")
var remoteCmdHost = flag.String("rch", "192.168.1.50:10001", "remote cmd host")
var localHost = flag.String("lh", "localhost:8000", "remote cmd host")

func main() {
	flag.Parse()

	if *mode == "server" {
		server.Listen(*port, *cmdPort)
	} else {
		client.Listen(*remoteCmdHost, *localHost)
	}

}
