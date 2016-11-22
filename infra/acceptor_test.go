package infra

import (
	"fmt"
	"testing"
)

func TestDefaults(t *testing.T) {

	s := extractConfig(nil)

	if s == nil {
		t.Fatal("No defaults returned")
	}

	fmt.Println((*s).host, (*s).port)
	if s.host != host_default+":" || s.port != port_default ||
		s.protocol != protocol_default {

		t.Fatal("Unexpected defaults returned")
	}

	var a Acceptor
	err := a.Open(nil)

	if err != nil {
		t.Fatal("Couldn't open acceptor", err)
	}

	err = a.Run()

	if err != nil {
		t.Fatal("Couldn't open acceptor", err)
	}

}
