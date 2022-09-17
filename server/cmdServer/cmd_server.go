package cmdServer

import (
	"awesomeProject2/server/svcContext"
	"bufio"
	"fmt"
	"net"
	"time"
)

type CMDHandler struct {
	SvcCtx  *svcContext.SVCContext
	CmdPort int
	conn    *net.Conn
}

func (h *CMDHandler) Start() {
	go h.cmdServer()

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", h.CmdPort))
	if err != nil {
		fmt.Println("Error listening", err.Error())
		return //终止程序
	}
	// 监听并接受来自客户端的连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return // 终止程序
		}
		go h.doCMDChannalStuff(conn)
	}
}

func (h *CMDHandler) doCMDChannalStuff(conn net.Conn) {

	bufConn := bufio.NewReader(conn)
	fun := []byte{0}
	if _, err := bufConn.Read(fun); err != nil {
		conn.Close()
		fmt.Printf("[ERR] socks: read : %v", err)
		return
	}
	//cmd
	if fun[0] == 0 {
		fmt.Printf("New cmd client connected!\n")
		h.SvcCtx.NewConnNotiCh <- conn

		for {
			conn.SetReadDeadline(time.Now().Add(10 * time.Second))
			if _, err := bufConn.Read(fun); err != nil {
				conn.Close()
				fmt.Printf("[ERR] socks: read2 : %v", err)
				return
			}
			fun := []byte{5}
			conn.Write(fun)
		}

	}
	// new data channal
	if fun[0] == 1 {
		fmt.Printf("New data client connected\n")
		h.SvcCtx.ConnCh <- conn
	}
}

func (h *CMDHandler) cmdServer() {
	for {
		select {
		case needChannal, _ := <-h.SvcCtx.CmdCh:
			fmt.Printf("need new channal: %v \n", needChannal)
			bs := []byte{1}
			fmt.Printf("send conn request to client %v\n", h.conn)
			if h.conn != nil {
				(*h.conn).Write(bs)
			}
			fmt.Printf("send to client to create a new channal\n")
		case conn, _ := <-h.SvcCtx.NewConnNotiCh:
			h.conn = &conn
		}
	}
}