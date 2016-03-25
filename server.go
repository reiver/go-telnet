package telnet


import (
	"io"
	"net"
)


// ListenAndServe listens on the TCP network address `addr` and then spawns a call to the ServeTELNET
// method on the `handler` to serve each incoming connection.
//
// For a very simple example:
//
//	package main
//	
//	import (
//		"github.com/reiver/go-telnet"
//	)
//	
//	func main() {
//	
//		var handler telnet.Handler = telnet.EchoHandler
//	
//		err := telnet.ListenAndServe(":5555", handler)
//		if nil != err {
//			//@TODO: Handle this error better.
//			panic(err)
//		}
//	}
func ListenAndServe(addr string, handler Handler) error {
	server := &Server{Addr: addr, Handler: handler}
	return server.ListenAndServe()
}


// Serve accepts an incoming TELNET client connection on the net.Listener `listener`.
func Serve(listener net.Listener, handler Handler) error {

	server := &Server{Handler: handler}
	return server.Serve(listener)
}


// A Server defines parameters of a running TELNET server.
//
// For a simple example:
//
//	package main
//	
//	import (
//		"github.com/reiver/go-telnet"
//	)
//	
//	func main() {
//	
//		var handler telnet.Handler = telnet.EchoHandler
//	
//		server := &telnet.Server{
//			Addr:":5555",
//			Handler:handler,
//		}
//	
//		err := server.ListenAndServe()
//		if nil != err {
//			//@TODO: Handle this error better.
//			panic(err)
//		}
//	}
type Server struct {
	Addr    string  // TCP address to listen on; ":telnet" or ":telnets" if empty (when used with ListenAndServe or ListenAndServeTLS respectively).
	Handler Handler // handler to invoke; telnet.EchoServer if nil
}


// ListenAndServe listens on the TCP network address 'server.Addr' and then spawns a call to the ServeTELNET
// method on the 'server.Handler' to serve each incoming connection.
//
// For a simple example:
//
//	package main
//	
//	import (
//		"github.com/reiver/go-telnet"
//	)
//	
//	func main() {
//	
//		var handler telnet.Handler = telnet.EchoHandler
//	
//		server := &telnet.Server{
//			Addr:":5555",
//			Handler:handler,
//		}
//	
//		err := server.ListenAndServe()
//		if nil != err {
//			//@TODO: Handle this error better.
//			panic(err)
//		}
//	}
func (server *Server) ListenAndServe() error {

	addr := server.Addr
	if "" == addr {
		addr = ":telnet"
	}


	listener, err := net.Listen("tcp", addr)
	if nil != err {
		return err
	}


	return server.Serve(listener)
}


// Serve accepts an incoming TELNET client connection on the net.Listener `listener`.
func (server *Server) Serve(listener net.Listener) error {

	defer listener.Close()


	handler := server.Handler
	if nil == handler {
//@TODO: Should this be a "ShellHandler" instead, that gives a shell-like experience by default
//       If this is changd, then need to change the comment in the "type Server struct" definition.
		handler = EchoHandler
	}


	for {
		// Wait for a new TELNET client connection.
		conn, err := listener.Accept()
		if err != nil {
//@TODO: Could try to recover from certain kinds of errors. Maybe waiting a while before trying again.
			return err
		}

		// Handle the new TELNET client connection by spawning
		// a new goroutine.
		go func(c net.Conn) {
//@TODO: Add proper context.
			var ctx Context = nil

			var w io.Writer = c
			var r io.Reader = c

			handler.ServeTELNET(ctx, w, r)
			c.Close()
		}(conn)
	}
}
