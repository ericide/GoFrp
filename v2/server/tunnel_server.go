package server

import (
	"GoFrp/v2/constant"
	"GoFrp/v2/model"
	"GoFrp/v2/util"
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type TunnelServer struct {
	Port   int
	tunnel *Tunnel
}

func (s *TunnelServer) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", s.Port))
	if err != nil {
		log.Println("Error listening", err.Error())
		return //终止程序
	}
	// 监听并接受来自客户端的连接
	log.Println("start Listening tunnel slave conn", s.Port)
	for {
		conn, err := listener.Accept()
		log.Printf("Tunnel accepted new slave connection %v \n", conn)
		if err != nil {
			log.Println("Error accepting", err.Error())
			return // 终止程序
		}

		log.Println("New slave connect, start auth")
		err = s.auth(conn)
		if err != nil {
			log.Println("auth error", err.Error())
			continue
		}
		log.Println("auth success, start tunnel work")

		tunnel := NewTunnel(conn)
		s.tunnel = tunnel
		go tunnel.Start()
	}
}

func (s *TunnelServer) auth(conn net.Conn) error {
	fun := []byte{0, 0, 0, 0, 0, 0, 0, 0}

	log.Println("Auth", conn)
	_, err := io.ReadAtLeast(conn, fun, len(fun))
	if err != nil {
		return err
	}

	b0 := fun[0]
	b1 := fun[1]
	b2 := fun[2]
	b3 := fun[3]
	b4 := fun[4]
	b5 := fun[5]
	b6 := fun[6]
	b7 := fun[7]

	// a simple authorized method for testing
	if b0 == 1 && b1 == 7 && b2 == 5 && b3 == 8 && b4 == 6 && b5 == 5 && b6 == 9 && b7 == 0 {
		return nil
	}

	return errors.New("error conn")
}

func (s *TunnelServer) AddNewTunnelTask(conn *net.Conn) {
	s.tunnel.AddNewTunnelTask(conn)
}

type Tunnel struct {
	tunnelConn *net.Conn
	extConns   map[int64]*net.Conn
	index      int64
	writeCh    chan model.DataObject
}

func NewTunnel(conn net.Conn) *Tunnel {
	return &Tunnel{
		tunnelConn: &conn,
		extConns:   make(map[int64]*net.Conn),
		index:      0,
		writeCh:    make(chan model.DataObject),
	}
}

func (t *Tunnel) Close() {
	(*t.tunnelConn).Close()
}

func (t *Tunnel) Start() {
	go t.startTunnelWriteQueue()

	log.Println("Start liten to tunnel", *t.tunnelConn)
	for {
		fun := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		_, err := io.ReadAtLeast(*t.tunnelConn, fun, len(fun))

		if err != nil {
			log.Println(err)
			t.Close()
			return
		}

		identity, length, method, err := util.VerifyDataHeader(fun)
		if err != nil {
			t.Close()
			log.Println(err)
			return
		}

		switch method {
		case constant.MethodClose:
			t.closeExtConn(identity)
		case constant.MethodData:
			data := make([]byte, length)
			_, err := io.ReadAtLeast(*t.tunnelConn, data, len(data))
			if err != nil {
				log.Println(err)
				t.Close()
				return
			}
			t.transferData(identity, &data)
		case constant.MethodPing:
			log.Println("received heart ping:", identity)
			t.responseToPing(identity)
		case constant.MethodPong:

		}
	}
}

func (t *Tunnel) closeExtConn(index int64) {
	if clientConn, ok := t.extConns[index]; ok {
		err := (*clientConn).Close()
		if err != nil {
			log.Println(err)
		}
	} else {
		fmt.Println(index, "not in map when close")
	}
}

func (t *Tunnel) transferData(index int64, data *[]byte) {
	if clientConn, ok := t.extConns[index]; ok {
		log.Println("Write to ext", index, "length", len(*data))
		_, err := (*clientConn).Write(*data)
		if err != nil {
			log.Println(err)
		}
	} else {
		fmt.Println(index, "not in map when write")
	}
}

func (t *Tunnel) AddNewTunnelTask(conn *net.Conn) {
	connIndex := t.index
	t.extConns[connIndex] = conn
	t.index += 1

	go t.listenToExternalConn(connIndex, conn)

}

func (t *Tunnel) listenToExternalConn(index int64, conn *net.Conn) {
	bufConn := bufio.NewReader(*conn)
	for {
		datas := make([]byte, 1024)
		length, err := bufConn.Read(datas)
		if err != nil {
			fmt.Println("Ext client err: ", index, err)
			t.closeSideConn(index)
			return
		}

		header := util.CreateDataHeader(index, length, constant.MethodData)
		dataObj := model.DataObject{
			Pre:        header,
			Data:       &datas,
			DataLength: int64(length),
		}
		t.writeCh <- dataObj
	}
}
func (t *Tunnel) responseToPing(index int64) {
	header := util.CreateDataHeader(index, 0, constant.MethodPong)
	dataObj := model.DataObject{
		Pre:        header,
		Data:       nil,
		DataLength: 0,
	}
	t.writeCh <- dataObj
}
func (t *Tunnel) SendPing(index int64) {
	header := util.CreateDataHeader(index, 0, constant.MethodPing)
	dataObj := model.DataObject{
		Pre:        header,
		Data:       nil,
		DataLength: 0,
	}
	t.writeCh <- dataObj
}
func (t *Tunnel) closeSideConn(index int64) {
	header := util.CreateDataHeader(index, 0, constant.MethodClose)
	dataObj := model.DataObject{
		Pre:        header,
		Data:       nil,
		DataLength: 0,
	}
	t.writeCh <- dataObj
}

func (t *Tunnel) startTunnelWriteQueue() {
	for {
		select {
		case dataObj, _ := <-t.writeCh:

			log.Println("Write to tunnel:", dataObj.DataLength)

			_, err := (*t.tunnelConn).Write(*dataObj.Pre)
			if err != nil {
				log.Println("write header err", err)
				t.Close()
				return
			}
			if dataObj.DataLength != 0 {
				_, err = (*t.tunnelConn).Write((*dataObj.Data)[0:dataObj.DataLength])
				if err != nil {
					log.Println("write data err", err)
					t.Close()
					return
				}
			}
		}
	}
}
