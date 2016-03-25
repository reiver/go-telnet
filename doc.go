/*
Package telnet provides TELNET and TELNETS client and server implementations
in a style similar to the "net/http" library that is part of the Go standard library,
including support for "middleware"; TELNETS is secure TELNET, with the TELNET protocol
over a secured TLS (or SSL) connection.


Example TELNET Server

ListenAndServe starts a TELNET server with a given address and handler.

	handler := telnet.EchoHandler
	
	err := telnet.ListenAndServe(":23", handler)
	if nil != err {
		panic(err)
	}


Example TELNETS Server

ListenAndServeTLS starts a (secure) TELNETS server with a given address and handler,
using the specified "cert.pem" and "key.pem" files.

	handler := telnet.EchoHandler
	
	err := telnet.ListenAndServeTLS(":992", "cert.pem", "key.pem", handler)
	if nil != err {
		panic(err)
	}


TELNET vs TELNETS

If you are communicating over the open Internet, you should be using (the secure) TELNETS and ListenAndServeTLS.

If you are communicating just on localhost, then using just (the un-secure) TELNET and telnet.ListenAndServe may be OK.

If you are not sure which to use, use TELNETS and ListenAndServeTLS.


Example TELNET Shell Server

The previous 2 exaple servers were very very simple. Specifically, they just echoed back whatever
you submitted to it.

If you typed:

	Apple Banana Cherry\r\n

... it would send back:

	Apple Banana Cherry\r\n

(Exactly the same data you sent it.)

A more useful TELNET server can be made using the "github.com/reiver/go-telnet/telsh" sub-package.

The `telsh` sub-package provides "middleware" that enables you to create a "shell" interface (also
called a "command line interface" or "CLI") which most people would expect when using TELNET OR TELNETS.

For example:


	package main
	
	
	import (
		"github.com/reiver/go-oi"
		"github.com/reiver/go-telnet"
		"github.com/reiver/go-telnet/telsh"
		
		"time"
	)
	
	
	func main() {
		
		shellHandler := telsh.NewShellHandler()

		commandName := "date"
		shellHandler.Register(commandName, danceProducer)

		commandName = "animate"
		shellHandler.Register(commandName, animateProducer)
		
		addr := ":23"
		if err := telnet.ListenAndServe(addr, shellHandler); nil != err {
			panic(err)
		}
	}

Note that in the example, so far, we have registered 2 commands: `date` and `animate`.

For this to actually work, we need to have code for the `date` and `animate` commands.

The actual implemenation for the `date` command could be done like the following:

	func dateHandlerFunc(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser) error {
		const layout = "Mon Jan 2 15:04:05 -0700 MST 2006"
		s := time.Now().Format(layout)
		
		if _, err := oi.LongWriteString(stdout, s); nil != err {
			return err
		}
		
		return nil
	}
	
	
	func dateProducerFunc(ctx telsh.Context, name string, args ...string) telsh.Handler{
		return telsh.PromoteHandlerFunc(dateHandler)
	}
	
	
	var dateProducer = ProducerFunc(dateProducerFunc)

Note that your "real" work is in the `dateHandlerFunc` func.

And the actual implementation for the `animate` command could as done as follows:

	func animateHandlerFunc(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser) error {
		
		for i:=0; i<20; i++ {
			oi.LongWriteString(stdout, "\r⠋")
			time.Sleep(50*time.Millisecond)
			
			oi.LongWriteString(stdout, "\r⠙")
			time.Sleep(50*time.Millisecond)
			
			oi.LongWriteString(stdout, "\r⠹")
			time.Sleep(50*time.Millisecond)
			
			oi.LongWriteString(stdout, "\r⠸")
			time.Sleep(50*time.Millisecond)
			
			oi.LongWriteString(stdout, "\r⠼")
			time.Sleep(50*time.Millisecond)
			
			oi.LongWriteString(stdout, "\r⠴")
			time.Sleep(50*time.Millisecond)
			
			oi.LongWriteString(stdout, "\r⠦")
			time.Sleep(50*time.Millisecond)
			
			oi.LongWriteString(stdout, "\r⠧")
			time.Sleep(50*time.Millisecond)
			
			oi.LongWriteString(stdout, "\r⠇")
			time.Sleep(50*time.Millisecond)
			
			oi.LongWriteString(stdout, "\r⠏")
			time.Sleep(50*time.Millisecond)
		}
		oi.LongWriteString(stdout, "\r \r\n")
		
		return nil
	}
	
	
	func animateProducerFunc(ctx telsh.Context, name string, args ...string) telsh.Handler{
		return telsh.PromoteHandlerFunc(animateHandler)
	}
	
	
	var animateProducer = ProducerFunc(animateProducerFunc)

Again, note that your "real" work is in the `animateHandlerFunc` func.

Generating PEM Files

If you are using the telnet.ListenAndServeTLS func or the telnet.Server.ListenAndServeTLS method, you will need to
supply "cert.pem" and "key.pem" files.

If you do not already have these files, the Go soure code contains a tool for generating these files for you.

It can be found at:

	$GOROOT/src/crypto/tls/generate_cert.go

So, for example, if your `$GOROOT` is the "/usr/local/go" directory, then it would be at:

	/usr/local/go/src/crypto/tls/generate_cert.go

If you run the command:

	go run $GOROOT/src/crypto/tls/generate_cert.go --help

... then you get the help information for "generate_cert.go".

Of course, you would replace or set `$GOROOT` with whatever your path actualy is. Again, for example,
if your `$GOROOT` is the "/usr/local/go" directory, then it would be:

	go run /usr/local/go/src/crypto/tls/generate_cert.go --help

To demonstrate the usage of "generate_cert.go", you might run the following to generate certificates
that were bound to the hosts `127.0.0.1` and `localhost`:

	go run /usr/local/go/src/crypto/tls/generate_cert.go --ca --host='127.0.0.1,localhost'


If you are not sure where "generate_cert.go" is on your computer, on Linux and Unix based systems, you might
be able to find the file with the command:

	locate /src/crypto/tls/generate_cert.go

(If it finds it, it should output the full path to this file.)


*/
package telnet
