package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:5000")
	if err != nil {
		fmt.Println("client err=", err)
		return
	}
	defer conn.Close() // 关闭连接
	//客户端可以发送单行数据


	_, err = conn.Write([]byte{0})
	if err != nil {
		fmt.Println("conn.write err=", err)
	}

	reader := bufio.NewReader(conn)
	for {

		fun := []byte{0}
		if _, err := reader.Read(fun); err != nil {
			fmt.Printf("[ERR] socks: Failed to get version byte: %v", err)
			return
		}
		go createNewConn()

	}
	//fmt.Printf("客户端发送了%d 字节的数量", n)

}

func createNewConn() {
	conn, err := net.Dial("tcp", "0.0.0.0:5000")
	if err != nil {
		fmt.Println("client err=", err)
		return
	}
	conn.Write([]byte{1})

	conn2, err := net.Dial("tcp", "192.168.1.10:8082")
	if err != nil {
		fmt.Println("client err=", err)
		return
	}

	fmt.Println("start transmit")

	errCh := make(chan error, 2)
	go proxy2(conn, conn2, errCh)
	go proxy2(conn2, conn, errCh)
	for i := 0; i < 2; i++ {
		e := <-errCh
		if e != nil {
			// return from this function closes target (and conn).
			return
		}
	}
}

func proxy2(dst io.Writer, src io.Reader, errCh chan error) {
	_, err := io.Copy(dst, src)
	errCh <- err
}

