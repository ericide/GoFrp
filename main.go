package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var cmd_ch = make(chan int)
var ch2 = make(chan net.Conn)

func main() {
	go doCMDChannalServer(5000)
	go doListenServer(5001)

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
}


func doCMDChannalServer(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
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
		go doCMDChannalStuff(conn)
	}
}

func doCMDChannalStuff(conn net.Conn) {

	bufConn := bufio.NewReader(conn)
	fun := []byte{0}
	if _, err := bufConn.Read(fun); err != nil {
		fmt.Printf("[ERR] socks: Failed to get version byte: %v", err)
		return
	}
	//cmd
	if fun[0] == 0 {
		fmt.Printf("New cmd client connected!\n")
		go cmdServer(conn)
	}
	// new data channal
	if fun[0] == 1 {
		fmt.Printf("New data client connected\n")
		ch2<- conn
	}

}

func cmdServer(conn net.Conn) {
	for {
		needChannal, _ := <-cmd_ch
		fmt.Printf("need new channal: %v \n", needChannal)
		bs := []byte{0}
		conn.Write(bs)
		fmt.Printf("send to client to create a new channal\n")
	}
}


func doListenServer(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
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
		go doListenStuff(conn)
	}
}

func doListenStuff(conn net.Conn) {

	fmt.Printf("New external request received!\n")

	cmd_ch<- 1
	conn2, _ := <-ch2

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