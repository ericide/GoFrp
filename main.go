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
	"os"
	"sync"
)

var configPath = flag.String("c", "", "config file path")

func main() {

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	defaultConfig := path.Join(exPath, "config.json")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.Parse()

	if "" == *configPath {
		configPath = &defaultConfig
	}

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

	<-ctx.Done()
	fmt.Println("bye")
}
