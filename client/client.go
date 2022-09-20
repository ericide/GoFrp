package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func connectRemoteController(remoteCmdHost string, localHost string) {
	conn, err := net.Dial("tcp", remoteCmdHost)
	if err != nil {
		fmt.Println("client err=", err)
		return
	}
	defer conn.Close() // 关闭连接
	fmt.Println("conn successful")
	_, err = conn.Write([]byte{0})
	if err != nil {
		fmt.Println("conn.write err=", err)
	}

	go heartRate(conn)

	reader := bufio.NewReader(conn)
	for {

		fun := []byte{0}
		if _, err := reader.Read(fun); err != nil {
			log.Printf("[ERR] socks: Failed to get version byte: %v", err)
			return
		}

		if fun[0] == 1 {
			log.Printf("new connect request %v \n ", fun)
			go createNewConn(remoteCmdHost, localHost)
		} else {
			log.Printf("hert rate response: %v \n", fun[0])
		}
	}
}
func heartRate(conn net.Conn) {
	for {
		time.Sleep(5 * time.Second)
		_, err := conn.Write([]byte{0})
		if err != nil {
			return
		}
	}
}

func Listen(remoteCmdHost string, localHost string) {

	for {
		connectRemoteController(remoteCmdHost, localHost)
		time.Sleep(2 * time.Second)
	}
	//客户端可以发送单行数据

	//fmt.Printf("客户端发送了%d 字节的数量", n)

}

func createNewConn(remoteCmdHost string, localHost string) {
	conn, err := net.Dial("tcp", remoteCmdHost)
	if err != nil {
		fmt.Println("client err=", err)
		return
	}
	conn.Write([]byte{1})

	conn2, err := net.Dial("tcp", localHost)
	if err != nil {
		fmt.Println("client err=", err)
		return
	}

	fmt.Println("start transmit")

	errCh := make(chan error, 2)
	go proxy2("local -> remote", conn, conn2, errCh)
	go proxy2("remote -> local", conn2, conn, errCh)
	for i := 0; i < 2; i++ {
		e := <-errCh
		if e != nil {
			// return from this function closes target (and conn).
			return
		}
	}
}

func proxy2(des string, dst io.Writer, src io.Reader, errCh chan error) {
	num, err := io.Copy(dst, src)
	log.Printf("num: %v, des: %s err: %v direction: %v -> %v", num, des, err, src, dst)
	errCh <- err
}
