package infra

import (
	"net"
)

// Type declaration for concurency model
type ConcurrencyModel int

// Different concurrency models at work
const (
	RoPC   ConcurrencyModel = iota // Go-Routine Per Connection Model
	RoPool                         // Routine Pool
)

// Service handler registered by the application to do real work
type ServiceHandler interface {
	HandleEvents(c net.Conn) error
}
