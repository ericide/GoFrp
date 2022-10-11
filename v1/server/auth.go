package server

import (
	"errors"
	"io"
	"log"
	"net"
)

func auth(conn net.Conn) error {
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
