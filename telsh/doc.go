/*
Package telsh provides "middle" (for the telnet package) that can be used to implement a TELNET server
that provides a "shell" interface.

Shell interfaces you may be familiar with include: "bash", "csh", "sh", "zsk", etc.

TELNET Server

Here is an example usage:

	package main
	
	import (
		"github.com/reiver/go-oi"
		"github.com/reiver/go-telnet"
		"github.com/reiver/go-telnet/telsh"

		"io"
	)

	func main() {
		
		telnetHandler := telsh.NewShellHandler()
		
		if err := telnetHandler.RegisterElse(
			telsh.ProducerFunc(
				func(ctx telsh.Context, name string, args ...string) telsh.Handler {
					return telsh.PromoteHandlerFunc(
						func(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser) error {
							oi.LongWrite(stdout, []byte{'w','a','t','?', '\r','\n'})
							
							return nil
						},
					)
				},
			),
		); nil != err {
			panic(err)
		}
		
		if err := telnetHandler.Register("help",
			telsh.ProducerFunc(
				func(ctx telsh.Context, name string, args ...string) telsh.Handler {
				return telsh.PromoteHandlerFunc(
						func(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser) error {
							oi.LongWrite(stdout, []byte{'r','t','f','m','!', '\r','\n'})
							
							return nil
						},
					)
				},
			),
		); nil != err {
			panic(err)
		}
		
		err := telnet.ListenAndServe(":5555", telnetHandler)
		if nil != err {
			//@TODO: Handle this error better.
			panic(err)
		}
	}

*/
package telsh
