package svcContext

import "net"

type SVCContext struct {
	CmdCh         chan int
	ConnCh        chan net.Conn
	NewConnNotiCh chan net.Conn
}
