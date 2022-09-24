package server

import (
	"GoFrp/v1/server/cmdServer"
	"GoFrp/v1/server/svcContext"
	"fmt"
	"io"
	"log"
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
		log.Println("Error listening", err.Error())
		return //终止程序
	}
	// 监听并接受来自客户端的连接
	for {
		log.Println("start Accept", port)

		conn, err := listener.Accept()
		log.Printf("Accepted %v \n", conn)

		if err != nil {
			log.Println("Error accepting", err.Error())
			return // 终止程序
		}
		go doListenStuff(ctx, conn)
	}
}

func doListenStuff(ctx *svcContext.SVCContext, conn net.Conn) {

	log.Printf("New external request received!\n")

	ctx.CmdCh <- 1
	connFromFrpClient, _ := <-ctx.ConnCh

	log.Printf("start transmit data!\n")

	errCh := make(chan error, 2)
	go proxy("frp client -> real client", conn, connFromFrpClient, errCh)
	go proxy("real client -> frp client", connFromFrpClient, conn, errCh)

	<-errCh

	err := conn.Close()
	log.Printf("close err 1 %v", err)
	err = connFromFrpClient.Close()
	log.Printf("close err 2 %v", err)

	<-errCh

	log.Printf("close %v, %v", conn, connFromFrpClient)

}

func proxy(des string, dst io.Writer, src io.Reader, errCh chan error) {
	num, err := io.Copy(dst, src)
	log.Printf("num: %v, des: %s err: %v direction: %v -> %v", num, des, err, src, dst)
	errCh <- err
}
