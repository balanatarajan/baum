package infra

import (
	"bufio"
	"fmt"
	"net"
	"testing"
	"time"
)

type myConnSH struct {
	cnt int // Number of times we are called
}

func (m *myConnSH) HandleEvents(c net.Conn) error {

	fmt.Println("Connection details ", c)

	w := bufio.NewWriter(c)

	for i := 0; i != 10; i++ {

		n, err := w.Write([]byte("What a great day!"))

		if err != nil {
			fmt.Println("What is the error ", err)
		}

		err = w.Flush()

		if err != nil {
			fmt.Println("What is the flush error ", err)
		}

		fmt.Println("Bytes written", n)

		time.Sleep(100 * time.Millisecond)
	}

	w.Write([]byte("Done"))
	w.Flush()
	return nil
}

func TestConnDefaults(t *testing.T) {

	s := extractConfig(nil)

	if s == nil {
		t.Fatal("No defaults returned")
	}

	fmt.Println((*s).host, (*s).port)
	if s.host != host_default+":" || s.port != port_default ||
		s.protocol != protocol_default {

		t.Fatal("Unexpected defaults returned")
	}

	var sh myConnSH

	c := NewConnector(RoPC, &sh)

	c.Connect(nil)

	c.Wait()
	c.Close()

}
