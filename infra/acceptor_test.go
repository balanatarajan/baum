package infra

import (
	"bufio"
	"fmt"
	"net"
	"testing"
)

type myAccSH struct {
	cnt int // Number of times we are called
}

func (m *myAccSH) HandleEvents(c net.Conn) error {
	(*m).cnt++

	nr := bufio.NewReader(c)
	fmt.Println(nr.Buffered())

	for {
		op := make([]byte, 5)

		n, err := nr.Read(op)

		if n == 4 || err != nil {
			break
		}

		if n != 0 {
			fmt.Println("Got ", string(op[:]))
		}

	}

	return nil
}

var sh myAccSH

func TestAccDefaults(t *testing.T) {

	s := extractConfig(nil)

	if s == nil {
		t.Fatal("No defaults returned")
	}

	fmt.Println((*s).host, (*s).port)
	if s.host != host_default+":" || s.port != port_default ||
		s.protocol != protocol_default {

		t.Fatal("Unexpected defaults returned")
	}

	a := Acceptor{Sh: &sh, Cm: RoPC}

	err := a.Open(nil)
	defer a.Close()

	if err != nil {
		t.Fatal("Couldn't open acceptor", err)
	}

	err = a.Run()

	if err != nil {
		t.Fatal("Couldn't open acceptor", err)
	}

}
