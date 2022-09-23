package slave

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type DataObject struct {
	Pre        *[]byte
	Data       *[]byte
	DataLength int64
}

var ConnCh = make(chan DataObject)
var maps = make(map[int64]net.Conn)

func Start(remoteCmdHost string, localHost string) {

	for {
		connectRemoteController(remoteCmdHost, localHost)
		time.Sleep(2 * time.Second)
	}
}

func serveForCOonn(conn net.Conn) {
	for {
		select {
		case dataObj, _ := <-ConnCh:
			conn.Write(*dataObj.Pre)
			conn.Write((*dataObj.Data)[0:dataObj.DataLength])
		}
	}
}

func connectRemoteController(remoteCmdHost string, localHost string) {
	conn, err := net.Dial("tcp", remoteCmdHost)
	if err != nil {
		fmt.Println("client err=", err)
		return
	}
	defer conn.Close() // 关闭连接
	fmt.Println("conn successful")

	go serveForCOonn(conn)

	for {

		fun := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		_, err := io.ReadAtLeast(conn, fun, len(fun))
		if err != nil {
			log.Println(err)
			return
		}
		var identy int = int(fun[3])*256*256*256 + int(fun[2])*256*256 + int(fun[1])*256 + int(fun[0])
		var length int = int(fun[7])*256*256*256 + int(fun[6])*256*256 + int(fun[5])*256 + int(fun[4])

		datas := make([]byte, length)
		_, err = io.ReadAtLeast(conn, datas, len(datas))
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println("will write", identy, length)

		if clientConn, ok := maps[int64(identy)]; ok {
			_, err = clientConn.Write(datas)
			if err != nil {
				log.Println(err)
				continue
			}
		} else {
			fmt.Println("new conn to", localHost)
			clientConn, err := net.Dial("tcp", localHost)
			if err != nil {
				fmt.Println("client err=", err)
				continue
			}
			maps[int64(identy)] = clientConn

			_, err = clientConn.Write(datas)
			if err != nil {
				log.Println(err)
				continue
			}
			go listenToResponse(int64(identy), clientConn)
		}
	}
}
func listenToResponse(index int64, conn net.Conn) {
	bufConn := bufio.NewReader(conn)
	for {
		datas := make([]byte, 512)
		n, err := bufConn.Read(datas)
		if err != nil {
			fmt.Println("read from", index, err)
			break
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

		dataObj := DataObject{
			Pre:        &pre,
			Data:       &datas,
			DataLength: int64(n),
		}

		ConnCh <- dataObj
	}
}
