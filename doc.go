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

ListenAndServeTLS starts a TELNETS server with a given address and handler.

	handler := telnet.EchoHandler
	
	err := telnet.ListenAndServeTLS(":992", "cert.pem", "key.pem", handler)
	if nil != err {
		panic(err)
	}

TELNET vs TELNETS

If you are communicating over the open Internet, you should be using TELNETS and ListenAndServeTLS.

If you are communicating just on localhost, then using just TELNET and telnet.ListenAndServe may be OK.

If you are not sure which to use, use TELNETS and ListenAndServeTLS.

*/
package telnet
