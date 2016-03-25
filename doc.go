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

ListenAndServeTLS starts a (secure) TELNETS server with a given address and handler.

	handler := telnet.EchoHandler
	
	err := telnet.ListenAndServeTLS(":992", "cert.pem", "key.pem", handler)
	if nil != err {
		panic(err)
	}


TELNET vs TELNETS

If you are communicating over the open Internet, you should be using (the secure) TELNETS and ListenAndServeTLS.

If you are communicating just on localhost, then using just (the un-secure) TELNET and telnet.ListenAndServe may be OK.

If you are not sure which to use, use TELNETS and ListenAndServeTLS.


Generating "cert.pem" and "key.pem" Files

If you are using the telnet.ListenAndServeTLS func or the telnet.Server.ListenAndServeTLS method, you will need to
supply "cert.pem" and "key.pem" files.

The Go soure code contains a tool for generating these files for you. It can be found at:

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


Finding "generate_cert.go"

If you are not sure where "generate_cert.go" is on your computer, on Linux and Unix based systems, you might
be able to find the file with the command:

	locate /src/crypto/tls/generate_cert.go

(If it finds it, it should output the full path to this file.)


*/
package telnet
