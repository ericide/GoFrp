package slave

import (
	"GoFrp/v1/constant"
	"GoFrp/v1/util"
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

	_, err = conn.Write([]byte{constant.LInkTypeSignal})
	if err != nil {
		fmt.Println("conn.write err=", err)
	}

	err = auth(conn)
	if err != nil {
		fmt.Println("conn.write err=", err)
	}

	go heartRate(conn)
	fmt.Println("start ")
	for {
		fun := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		_, err := io.ReadAtLeast(conn, fun, len(fun))

		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}

		identity, _, method, err := util.VerifyDataHeader(fun)

		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}

		switch method {
		case constant.MethodPong:
			break
		case constant.MethodApplyNewDataChannel:
			go createNewConn(remoteCmdHost, localHost, identity)
			break
		}
	}
}

func auth(conn net.Conn) error {
	fun := []byte{1, 7, 5, 8, 6, 5, 9, 0}
	conn.Write(fun)
	return nil
}

func heartRate(conn net.Conn) {
	var index int64 = 0
	for {
		time.Sleep(5 * time.Second)
		header := util.CreateDataHeader(index, 0, constant.MethodPing)
		_, err := conn.Write(*header)
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

func createNewConn(remoteCmdHost string, localHost string, index int64) {
	tunnel, err := net.Dial("tcp", remoteCmdHost)
	if err != nil {
		fmt.Println("client err=", err)
		return
	}

	_, err = tunnel.Write([]byte{constant.LInkTypeDataTunnel})
	if err != nil {
		fmt.Println("client err=", err)
		return
	}
	fmt.Println("replay to apple data tunnel")

	header := util.CreateDataHeader(index, 0, constant.MethodApplyNewDataChannel)
	_, err = tunnel.Write(*header)
	if err != nil {
		fmt.Println("client err=", err)
		return
	}

	conn, err := net.Dial("tcp", localHost)
	if err != nil {
		fmt.Println("client err=", err)
		return
	}

	fmt.Println("start transmit")

	errCh := make(chan error, 2)
	go proxy2("local -> remote", conn, tunnel, errCh)
	go proxy2("remote -> local", tunnel, conn, errCh)
	<-errCh
	conn.Close()
}

func proxy2(des string, dst io.Writer, src io.Reader, errCh chan error) {

	//for {
	//	bytes, err := io.ReadAll(src)
	//	if err != nil {
	//		errCh <- err
	//		return
	//	}
	//	dst.Write(bytes)
	//	fmt.Println(des, len(bytes))
	//}

	num, err := io.Copy(dst, src)
	log.Printf("num: %v, des: %s err: %v direction: %v -> %v", num, des, err, src, dst)
	errCh <- err
}
