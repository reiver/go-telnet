package telsh


import (
	"github.com/reiver/go-telnet"
)


// A Producer provides a Produce method which creates a Handler.
//
// Producer is an abstraction that represents a shell "command".
//
// Contrast this with a Handler, which is is an abstraction that
// represents a "running" shell "command".
//
// To use a metaphor, the differences between a Producer and a Handler,
// is like the difference between a program executable and actually running
// the program executable.
type Producer interface {
	Produce(telnet.Context, string, ...string) Handler
}


// ProducerFunc is an adaptor, that can be used to turn a func with the
// signature:
//
//	func(telnet.Context, string, ...string) Handler
//
// Into a Producer
type ProducerFunc func(telnet.Context, string, ...string) Handler


// Produce makes ProducerFunc fit the Producer interface.
func (fn ProducerFunc) Produce(ctx telnet.Context, name string, args ...string) Handler {
	return fn(ctx, name, args...)
}
