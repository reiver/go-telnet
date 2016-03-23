/*
Package telnet provides TELNET client and server implementations
in a style similar to the "net/http" library that is part of the Go standard library,
including support for "middleware".

TELNET Server

ListenAndServe starts a TELNET server with a given address and handler.

	handler := telnet.EchoHandler
	
	err := telnet.ListenAndServe(":5555", handler)
	if nil != err {
		panic(err)
	}

*/
package telnet
