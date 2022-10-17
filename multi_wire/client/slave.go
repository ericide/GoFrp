package client

import (
	"GoFrp/multi_wire/constant"
	"GoFrp/multi_wire/svcContext"
	"GoFrp/multi_wire/util"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func connectRemoteController(ctx *svcContext.SVCContext) {

	serverAddr := fmt.Sprintf("%s:%d", ctx.ServerHost, ctx.ServerPort)
	bindAddr := fmt.Sprintf("%s:%d", ctx.BindHost, ctx.BindPort)

	conn, err := net.Dial("tcp", serverAddr)
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

	err = auth(ctx, conn)
	if err != nil {
		fmt.Println("conn.write err=", err)
	}

	go heartRate(ctx, conn)
	fmt.Println("start ")
	for {
		identity, method, _, err := util.ReadDataPackage(ctx.Password, conn)

		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}

		switch method {
		case constant.MethodPong:
			break
		case constant.MethodApplyNewDataChannel:
			go createNewConn(ctx, serverAddr, bindAddr, identity)
			break
		}
	}
}

func auth(ctx *svcContext.SVCContext, conn net.Conn) error {
	body := []byte{1, 7, 5, 8, 6, 5, 9, 0}
	data := util.CreateDataPackage(ctx.Password, 0, constant.MethodSignalAuth, body)
	conn.Write(*data)
	return nil
}

func heartRate(ctx *svcContext.SVCContext, conn net.Conn) {
	var index int64 = 0
	for {
		time.Sleep(5 * time.Second)
		index++
		header := util.CreateDataPackage(ctx.Password, index, constant.MethodPing, nil)
		_, err := conn.Write(*header)
		if err != nil {
			return
		}
	}
}

func Start(ctx *svcContext.SVCContext) {

	for {
		connectRemoteController(ctx)
		time.Sleep(2 * time.Second)
	}
	//客户端可以发送单行数据

	//fmt.Printf("客户端发送了%d 字节的数量", n)

}

func createNewConn(ctx *svcContext.SVCContext, remoteCmdHost string, localHost string, index int64) {
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

	pack := util.CreateDataPackage(ctx.Password, index, constant.MethodApplyNewDataChannel, nil)
	_, err = tunnel.Write(*pack)
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
