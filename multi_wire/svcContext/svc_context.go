package svcContext

import (
	"sync"
)

type SVCContext struct {
	ApplyNewDataTunChan chan int64
	TaskMap             sync.Map
	ServerPort          int
	ServerHost          string
	BindPort            int
	BindHost            string
	Password            []byte
}
