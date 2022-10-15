package server

import (
	"GoFrp/v1/constant"
	"GoFrp/v1/svcContext"
	"GoFrp/v1/util"
	"fmt"
	"io"
	"log"
	"net"
)

func Start(ctx *svcContext.SVCContext) {

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", ctx.ServerPort))
	if err != nil {
		log.Println("Error listening", err.Error())
		return //终止程序
	}
	// 监听并接受来自客户端的连接
	var index int64 = 0
	log.Println("start Accept", ctx.ServerPort)
	for {
		index++
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		fmt.Println("new link come in", index)
		go doDistribute(conn, index, ctx)

	}
}

func doDistribute(conn net.Conn, index int64, ctx *svcContext.SVCContext) {
	fun := []byte{0}

	l, err := io.ReadAtLeast(conn, fun, len(fun))
	if err != nil {
		fmt.Println("new link err:", err)
		return
	}

	fmt.Println("new link first complete", index, l)

	if fun[0] == constant.LInkTypeSignal {
		go doSignalTunnel(conn, fun, ctx)
	} else if fun[0] == constant.LInkTypeDataTunnel {
		go doDataTunnel(conn, fun, ctx)
	} else { // data channal
		go doRequest(conn, index, fun, ctx)
	}
}

func doSignalTunnel(conn net.Conn, first []byte, ctx *svcContext.SVCContext) {

	fmt.Println("New Signal Link")

	err := auth(ctx, conn)
	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}

	//create new chan, each signal Tunnel has its own chan
	ctx.ApplyNewDataTunChan = make(chan int64)

	go runSignalWrite(conn, ctx)
	fmt.Println("start signal tunnel")
	for {
		identity, method, _, err := util.ReadDataPackage(ctx.Password, conn)
		if err != nil {
			conn.Close()
			return
		}
		// fmt.Println(body)

		if method == constant.MethodPing {
			pack := util.CreateDataPackage(ctx.Password, identity, constant.MethodPong, nil)
			conn.Write(*pack)
		}

	}
}
func runSignalWrite(conn net.Conn, ctx *svcContext.SVCContext) {
	for {
		identity := <-ctx.ApplyNewDataTunChan
		fmt.Println("apply link to remote slave")
		bytes := util.CreateDataPackage(ctx.Password, identity, constant.MethodApplyNewDataChannel, nil)
		_, err := conn.Write(*bytes)
		if err != nil {
			return
		}
	}
}

func doDataTunnel(conn net.Conn, first []byte, ctx *svcContext.SVCContext) {

	fmt.Println("New Data Link")

	identity, _, _, err := util.ReadDataPackage(ctx.Password, conn)

	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	value, ok := ctx.TaskMap.Load(identity)
	if ok != true {
		conn.Close()
		fmt.Println("map load not ok")
		return
	}
	if desChan, ok := value.(chan net.Conn); ok == true {
		desChan <- conn
	}
	ctx.TaskMap.Delete(identity)
}

func doRequest(conn net.Conn, index int64, first []byte, ctx *svcContext.SVCContext) {

	fmt.Println("New Request Link")

	myChan := make(chan net.Conn)

	ctx.ApplyNewDataTunChan <- index

	ctx.TaskMap.Store(index, myChan)

	tunnel, _ := <-myChan

	fmt.Println("Request get data tunnel")

	tunnel.Write(first)

	errCh := make(chan error, 2)

	fmt.Println("start tranmit")

	go proxy("<=", tunnel, conn, errCh)
	go proxy("=>", conn, tunnel, errCh)

	<-errCh

	conn.Close()
	tunnel.Close()
}

func proxy(des string, dst io.Writer, src io.Reader, errCh chan error) {

	//for {
	//	bytes, err := io.ReadAll(src)
	//	if err != nil {
	//		errCh <- err
	//		fmt.Println(des, err)
	//		return
	//	}
	//	dst.Write(bytes)
	//	fmt.Println(des, len(bytes))
	//}

	num, err := io.Copy(dst, src)
	log.Printf("num: %v, des: %s err: %v direction: %v -> %v", num, des, err, src, dst)
	errCh <- err
}
