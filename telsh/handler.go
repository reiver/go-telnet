package telsh


import (
	"io"
)


// Hander is an abstraction that represents a "running" shell "command".
//
// Constract this with a Producer, which is is an abstraction that
// represents a shell "command".
//
// To use a metaphor, the differences between a Producer and a Handler,
// is like the difference between a program executable and actually running
// the program executable.
//
// Conceptually, anything that implements the Hander, and then has its Producer
// registered with ShellHandler.Register() will be available as a command.
//
// Note that Handler was intentionally made to be compatible with
// "os/exec", which is part of the Go standard library.
type Handler interface {
	Run() error

	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)
}


// HandlerFunc is useful to write inline Producers, and provides an alternative to
// creating something that implements Handler directly.
//
// For example:
//
//	shellHandler := telsh.NewShellHandler()
//	
//	shellHandler.Register("five", telsh.ProducerFunc(
//		
//		func(ctx Context, name string, args ...string) telsh.Handler{
//		
//			return telsh.PromoteHandlerFunc(
//				
//				func(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser) error {
//					oi.LongWrite(stdout, []byte{'5', '\r', '\n'})
//					
//					return nil
//				},
//			)
//		},
//	))
//
// Note that PromoteHandlerFunc is used to turn a HandlerFunc into a Handler.
type HandlerFunc func(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser)error


type internalPromotedHandlerFunc struct {
	err error
	fn HandlerFunc
	stdin  io.ReadCloser
	stdout io.WriteCloser
	stderr io.WriteCloser

	stdinPipe  io.WriteCloser
	stdoutPipe io.ReadCloser
	stderrPipe io.ReadCloser
}


// PromoteHandlerFunc turns a HandlerFunc into a Handler.
func PromoteHandlerFunc(fn HandlerFunc) Handler {
	stdin,      stdinPipe := io.Pipe()
	stdoutPipe, stdout    := io.Pipe()
	stderrPipe, stderr    := io.Pipe()

	handler := internalPromotedHandlerFunc{
		err:nil,

		fn:fn,

		stdin:stdin,
		stdout:stdout,
		stderr:stderr,

		stdinPipe:stdinPipe,
		stdoutPipe:stdoutPipe,
		stderrPipe:stderrPipe,
	}

	return &handler
}


func (handler *internalPromotedHandlerFunc) Run() error {
	if nil != handler.err {
		return handler.err
	}

	handler.err =  handler.fn(handler.stdin, handler.stdout, handler.stderr)

	handler.stdin.Close()
	handler.stdout.Close()
	handler.stderr.Close()

	return handler.err
}

func (handler *internalPromotedHandlerFunc) StdinPipe() (io.WriteCloser, error) {
	if nil != handler.err {
		return nil, handler.err
	}

	return handler.stdinPipe, nil
}

func (handler *internalPromotedHandlerFunc) StdoutPipe() (io.ReadCloser, error) {
	if nil != handler.err {
		return nil, handler.err
	}

	return handler.stdoutPipe, nil
}

func (handler *internalPromotedHandlerFunc) StderrPipe() (io.ReadCloser, error) {
	if nil != handler.err {
		return nil, handler.err
	}

	return handler.stderrPipe, nil
}
