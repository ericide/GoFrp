package server

import (
	"awesomeProject2/server/cmdServer"
	"awesomeProject2/server/svcContext"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func Listen(port int, cmdPort int) {

	ctx := &svcContext.SVCContext{
		CmdCh:         make(chan int),
		ConnCh:        make(chan net.Conn),
		NewConnNotiCh: make(chan net.Conn),
	}

	go doListenServer(ctx, port)

	commandServer := cmdServer.CMDHandler{
		SvcCtx:  ctx,
		CmdPort: cmdPort,
	}
	go commandServer.Start()

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
}

func doListenServer(ctx *svcContext.SVCContext, port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		fmt.Println("Error listening", err.Error())
		return //终止程序
	}
	// 监听并接受来自客户端的连接
	for {
		fmt.Println("start Accept", port)

		conn, err := listener.Accept()
		fmt.Println("Accepted")

		if err != nil {
			fmt.Println("Error accepting", err.Error())
			return // 终止程序
		}
		go doListenStuff(ctx, conn)
	}
}

func doListenStuff(ctx *svcContext.SVCContext, conn net.Conn) {

	fmt.Printf("New external request received!\n")

	ctx.CmdCh <- 1
	conn2, _ := <-ctx.ConnCh

	fmt.Printf("start transmit data!\n")

	errCh := make(chan error, 2)
	go proxy(conn, conn2, errCh)
	go proxy(conn2, conn, errCh)
	for i := 0; i < 2; i++ {
		e := <-errCh
		if e != nil {
			// return from this function closes target (and conn).
			return
		}
	}
}

func proxy(dst io.Writer, src io.Reader, errCh chan error) {
	_, err := io.Copy(dst, src)
	errCh <- err
}
