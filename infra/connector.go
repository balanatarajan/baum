package infra

import (
	"fmt"
	"github.com/balanatarajan/baum/config"
	"net"
	"time"
)

// Feel comfortable with the name Acceptor from Acceptor pattern from
// Dr. Schmidt et. al. We will wrap listener behind the acceptor interface
// @TODO: We have to prep this for multi-routine scenarions

type Connector struct {
	Cm   ConcurrencyModel
	Sh   ServiceHandler
	conn net.Conn
	ch   chan int
}

func NewConnector(cm ConcurrencyModel, h ServiceHandler) *Connector {
	c := make(chan int)
	return &Connector{Cm: cm, Sh: h, ch: c}
}

func (a *Connector) Connect(config *config.Config) error {

	sc := extractConfig(config)

	var err error

	if a.Cm == RoPC {
		go func() {
			a.conn, err = net.DialTimeout((*sc).protocol, (*sc).host+(*sc).port, time.Second)
			if err != nil {
				fmt.Println("Cannot connect ", err)
			}
			a.Sh.HandleEvents(a.conn)
			a.ch <- 0
		}()
	} else {
		fmt.Println("ConcurrencyModel not supported")
	}

	if err != nil {
		return err
	}

	return err
}

func (a *Connector) Wait() error {

	if a.Cm == RoPC {
		<-a.ch
	}

	return nil
}

func (a *Connector) Close() error {

	close(a.ch)

	// ConcurrencyModel specific code needs to go in
	return a.conn.Close()
}
