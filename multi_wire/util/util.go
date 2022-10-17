package util

import (
	"bytes"
	"fmt"
	"io"
)

// pack
// length message
//        data
//        header body

func CreateDataPackage(pwd []byte, index int64, method byte, body []byte) *[]byte {
	var data = make([]byte, 6)
	//index
	data[0] = byte(index & 0xFF)
	data[1] = byte((index >> 8) & 0xFF)
	data[2] = byte((index >> 16) & 0xFF)
	data[3] = byte((index >> 24) & 0xFF)
	data[4] = method
	//sum
	var csum byte = 0
	for i := 0; i < 5; i++ {
		csum += data[i]
	}
	data[5] = csum

	if body != nil {
		data = append(data, body...)
	}

	message := AesEncrypt(data, pwd)

	pack := make([]byte, len(message)+2)

	// fmt.Println("Write Package:", len(message))

	//length
	pack[0] = byte(len(message) & 0xFF)
	pack[1] = byte((len(message) >> 8) & 0xFF)
	copy(pack[2:], message)

	return &pack

}

func ReadDataPackage(pwd []byte, conn io.ReadCloser) (int64, byte, []byte, error) {
	fun := []byte{0, 0}
	_, err := io.ReadAtLeast(conn, fun, len(fun))
	// fmt.Println(fun)
	if err != nil {
		conn.Close()
		return 0, 0, nil, err
	}
	length := int(fun[1])*256 + int(fun[0])

	// fmt.Println("Read Package:", length)

	fun = make([]byte, length)
	_, err = io.ReadAtLeast(conn, fun, len(fun))
	if err != nil {
		conn.Close()
		return 0, 0, nil, err
	}
	// fmt.Println(fun)
	return ParseDataMessage(pwd, fun)
}

func ParseDataMessage(pwd []byte, message []byte) (int64, byte, []byte, error) {
	data, err := AesDecrypt(message, pwd)

	if err != nil {
		return 0, 0, nil, fmt.Errorf("decrypt error, maybe discorrect password")
	}

	if len(message) < 6 { //message header mimimum is 6
		return 0, 0, nil, fmt.Errorf("unvalible message header")
	}

	var identity int64 = int64(data[3])*256*256*256 + int64(data[2])*256*256 + int64(data[1])*256 + int64(data[0])

	var method byte = data[4]
	var sum byte = data[5]
	// [0] [1] [2] [3] [4]      [5]         [6...]
	// [identity     ] [method] [sum check] [body]
	var csum byte = 0
	for i := 0; i < 5; i++ {
		csum += data[i]
	}

	if sum != csum {
		return 0, 0, nil, fmt.Errorf("error sum")
	}

	return identity, method, data[6:], nil

}

func ParsePassword(p string) []byte {
	ori := []byte(p)
	if len(ori) > 32 {
		return ori[:32]
	}

	padding := 32 - len(ori)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ori, padtext...)
}
