package server

import (
	"GoFrp/multi_wire/constant"
	"GoFrp/multi_wire/svcContext"
	"GoFrp/multi_wire/util"
	"errors"
	"fmt"
	"net"
)

func auth(ctx *svcContext.SVCContext, conn net.Conn) error {
	fmt.Println("Auth start")
	_, method, body, err := util.ReadDataPackage(ctx.Password, conn)

	if err != nil {
		return err
	}

	if method != constant.MethodSignalAuth {
		errors.New("error auth method")
	}

	if len(body) != 8 {
		errors.New("auth body error")
	}

	// a simple authorized method for testing
	if body[0] == 1 && body[1] == 7 && body[2] == 5 && body[3] == 8 && body[4] == 6 && body[5] == 5 && body[6] == 9 && body[7] == 0 {
		fmt.Println("Auth Success")
		return nil
	}

	return errors.New("error auth")
}
