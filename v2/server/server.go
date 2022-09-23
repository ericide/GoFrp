package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

var Gconn *net.Conn = nil
var maps map[int64]net.Conn = make(map[int64]net.Conn)

type DataObject struct {
	Pre        *[]byte
	Data       *[]byte
	DataLength int64
}

var ConnCh = make(chan DataObject)

func ListenServer(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		log.Println("Error listening", err.Error())
		return //终止程序
	}
	// 监听并接受来自客户端的连接
	var index int64 = 0
	for {
		log.Println("start Accept", port)

		conn, err := listener.Accept()
		log.Printf("Accepted %v \n", conn)

		if err != nil {
			log.Println("Error accepting", err.Error())
			return // 终止程序
		}
		index++
		fmt.Println(index, "into map")
		maps[index] = conn
		go readfromclient(index, conn)
	}
}

func readfromclient(index int64, conn net.Conn) {
	bufConn := bufio.NewReader(conn)
	for {
		datas := make([]byte, 512)
		n, err := bufConn.Read(datas)
		if err != nil {
			fmt.Println("client ", index, err)
			return
		}

		pre := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		//index
		pre[0] = byte(index & 0xFF)
		pre[1] = byte((index >> 8) & 0xFF)
		pre[2] = byte((index >> 16) & 0xFF)
		pre[3] = byte((index >> 24) & 0xFF)
		//length
		pre[4] = byte(n & 0xFF)
		pre[5] = byte((n >> 8) & 0xFF)
		pre[6] = byte((n >> 16) & 0xFF)
		pre[7] = byte((n >> 24) & 0xFF)

		fmt.Println(index, "send", n)

		dataObj := DataObject{
			Pre:        &pre,
			Data:       &datas,
			DataLength: int64(n),
		}

		ConnCh <- dataObj
	}
}

func ListenTunnelServer(port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		log.Println("Error listening", err.Error())
		return //终止程序
	}
	// 监听并接受来自客户端的连接

	go tunnelServer()
	for {
		log.Println("start Accept", port)

		conn, err := listener.Accept()
		log.Printf("Accepted %v \n", conn)
		if err != nil {
			log.Println("Error accepting", err.Error())
			return // 终止程序
		}
		Gconn = &conn

		go readTunnelServer()

	}
}

func tunnelServer() {
	for {
		select {
		case dataObj, _ := <-ConnCh:

			if Gconn == nil {
				continue
			}
			
			(*Gconn).Write(*dataObj.Pre)
			(*Gconn).Write((*dataObj.Data)[0:dataObj.DataLength])
		}
	}

}
func readTunnelServer() {
	for {
		fun := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		_, err := io.ReadAtLeast(*Gconn, fun, len(fun))
		if err != nil {
			log.Println(err)
			(*Gconn).Close()
			return
		}
		var identy int = int(fun[3])*256*256*256 + int(fun[2])*256*256 + int(fun[1])*256 + int(fun[0])
		var length int = int(fun[7])*256*256*256 + int(fun[6])*256*256 + int(fun[5])*256 + int(fun[4])

		datas := make([]byte, length)
		_, err = io.ReadAtLeast(*Gconn, datas, len(datas))
		if err != nil {
			log.Println(err)
			(*Gconn).Close()
			return
		}

		if clientConn, ok := maps[int64(identy)]; ok {
			_, err = clientConn.Write(datas)
			if err != nil {
				log.Println(err)
			}
		} else {
			fmt.Println(identy, "not in map")
		}
	}
}
