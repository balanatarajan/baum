package infra

import (
	"fmt"
	"github.com/balanatarajan/baum/config"
	"net"
	"strings"
)

type srvConfig struct {
	protocol,
	host,
	port string
}

// Feel comfortable with the name Acceptor from Acceptor pattern from
// Dr. Schmidt et. al. We will wrap listener behind the acceptor interface
type Acceptor struct {
	Cm       ConcurrencyModel
	Sh       ServiceHandler
	listener net.Listener
}

const host_default = "127.0.0.1"
const port_default = "10007"
const protocol_default = "tcp"

var defSrvConfig *srvConfig = defaultConfig()

func defaultConfig() *srvConfig {
	return &srvConfig{host: host_default + ":", port: port_default,
		protocol: protocol_default}
}

func extractConfig(config *config.Config) *srvConfig {

	if config == nil {
		return defSrvConfig
	}

	var s srvConfig
	var err error
	s.host, err = config.Find("ip")

	if err != nil || len(s.host) == 0 {
		s.host = host_default
	}

	if strings.HasSuffix(s.host, ":") == false {
		s.host += ":"
	}

	s.port, err = config.Find("port")

	if err != nil || len(s.port) == 0 {
		s.port = port_default
	}

	s.protocol, err = config.Find("protocol")

	if err != nil || len(s.protocol) == 0 {
		s.protocol = protocol_default
	}

	return &s
}

func (a *Acceptor) Open(config *config.Config) error {

	sc := extractConfig(config)

	var err error

	a.listener, err = net.Listen((*sc).protocol, (*sc).host+(*sc).port)

	if err != nil {
		return err
	}

	return err
}

func (a *Acceptor) Run() error {

	for {
		conn, _ := a.listener.Accept()
		if a.Cm == RoPC {
			go a.Sh.HandleEvents(conn)
		} else {
			fmt.Println("ConcurrencyModel not set")
		}
	}
}

func (a *Acceptor) Close() error {

	// ConcurrencyModel specific code needs to go in
	return a.listener.Close()
}
