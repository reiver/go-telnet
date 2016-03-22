/*
Package telnet provides TELNET client and server implementations.

The API this package provide was intentionally designed to seem similar to the
the standard "net/http" Go library.

Including the ability to create "middleware".

TELNET Server

ListenAndServe starts a TELNET server with a given address and handler.

	handler := telnet.EchoHandler
	
	err := telnet.ListenAndServe(":5555", handler)
	if nil != err {
		panic(err)
	}

*/
package telnet
