package util

import "errors"

func CreateDataHeader(index int64, length int, method byte) *[]byte {
	var pre = make([]byte, 8)
	//index
	pre[0] = byte(index & 0xFF)
	pre[1] = byte((index >> 8) & 0xFF)
	pre[2] = byte((index >> 16) & 0xFF)
	pre[3] = byte((index >> 24) & 0xFF)
	//length
	pre[4] = byte(length & 0xFF)
	pre[5] = byte((length >> 8) & 0xFF)
	pre[6] = method

	var csum byte = 0
	for i := 0; i < 7; i++ {
		csum += pre[i]
	}
	pre[7] = csum

	return &pre

}

func VerifyDataHeader(fun []byte) (int64, int, byte, error) {
	var identity int64 = int64(fun[3])*256*256*256 + int64(fun[2])*256*256 + int64(fun[1])*256 + int64(fun[0])
	var length int = int(fun[5])*256 + int(fun[4])
	var method byte = fun[6]
	var sum byte = fun[7]

	var csum byte = 0
	for i := 0; i < 7; i++ {
		csum += fun[i]
	}

	if sum != csum {
		return 0, 0, 0, errors.New("not correct sum")
	}

	return identity, length, method, nil

}
