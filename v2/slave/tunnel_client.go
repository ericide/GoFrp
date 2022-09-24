package slave

import (
	"GoFrp/v2/constant"
	"GoFrp/v2/util"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

type TunnelClient struct {
	TunnelServerAddress string
	ExtAddress          string
	tunnel              *Tunnel
}

func (s *TunnelClient) Start() {
	fmt.Println("Start connect to tunnel server", s.TunnelServerAddress)
	conn, err := net.Dial("tcp", s.TunnelServerAddress)
	if err != nil {
		fmt.Println("client err=", err)
		return
	}
	fmt.Println("Connect success")

	fmt.Println("Start auth")
	err = s.auth(conn)
	if err != nil {
		log.Println("Error accepting", err.Error())
		return
	}
	fmt.Println("Auth success")

	fmt.Println("Start tunnel work")
	tunnel := NewTunnel(conn, s.ExtAddress)
	s.tunnel = tunnel
	tunnel.Start()
}

func (s *TunnelClient) auth(conn net.Conn) error {
	fun := []byte{1, 7, 5, 8, 6, 5, 9, 0}
	conn.Write(fun)
	return nil
}

type Tunnel struct {
	ExtAddress string
	tunnelConn *net.Conn
	extConns   map[int64]*net.Conn
	index      int64
	writeCh    chan DataObject
}

func NewTunnel(conn net.Conn, extAddress string) *Tunnel {
	return &Tunnel{
		ExtAddress: extAddress,
		tunnelConn: &conn,
		extConns:   make(map[int64]*net.Conn),
		writeCh:    make(chan DataObject),
	}
}

func (t *Tunnel) Close() {
	(*t.tunnelConn).Close()
}

func (t *Tunnel) Start() {
	go t.startTunnelWriteQueue()
	for {
		fun := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		_, err := io.ReadAtLeast(*t.tunnelConn, fun, len(fun))
		if err != nil {
			log.Println(err)
			t.Close()
			return
		}

		log.Println("Tunnel", *t.tunnelConn, "have data")

		identity, length, method, err := util.VerifyDataHeader(fun)
		if err != nil {
			log.Println(err)
			t.Close()
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
		log.Println(index, "not in map when close")
	}
}

func (t *Tunnel) transferData(index int64, data *[]byte) {
	var conn *net.Conn
	if clientConn, ok := t.extConns[index]; ok {
		conn = clientConn
	} else {
		log.Println("Create new conn", index, t.ExtAddress)
		clientConn, err := net.Dial("tcp", t.ExtAddress)
		if err != nil {
			log.Println("Create new err=", err)
			t.closeSideConn(index)
			return
		}
		t.AddNewTunnelTask(index, &clientConn)
		conn = &clientConn
	}

	log.Println("Write to ext:", index, len(*data))
	_, err := (*conn).Write(*data)
	if err != nil {
		t.closeSideConn(index)
		log.Println(err)
		return
	}
}

func (t *Tunnel) AddNewTunnelTask(index int64, conn *net.Conn) {
	t.extConns[index] = conn
	go t.listenToExternalConn(index, conn)
}

func (t *Tunnel) listenToExternalConn(index int64, conn *net.Conn) {
	bufConn := bufio.NewReader(*conn)
	for {
		datas := make([]byte, 1024)
		length, err := bufConn.Read(datas)
		if err != nil {
			fmt.Println("client ", index, err)
			t.closeSideConn(index)
			return
		}
		log.Println("Read from ext", index, "Length", length)
		header := util.CreateDataHeader(index, length, constant.MethodData)
		dataObj := DataObject{
			Pre:        header,
			Data:       &datas,
			DataLength: int64(length),
		}
		t.writeCh <- dataObj
	}
}

func (t *Tunnel) closeSideConn(index int64) {
	header := util.CreateDataHeader(index, 0, constant.MethodClose)
	dataObj := DataObject{
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

			log.Println("Write to tunnel:", dataObj.DataLength, *t.tunnelConn)

			_, err := (*t.tunnelConn).Write(*dataObj.Pre)
			if err != nil {
				log.Println(err)
				t.Close()
				return
			}
			if dataObj.DataLength != 0 {
				_, err = (*t.tunnelConn).Write((*dataObj.Data)[0:dataObj.DataLength])
				if err != nil {
					log.Println(err)
					t.Close()
					return
				}
			}
		}
	}
}
