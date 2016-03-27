package telnet


import (
	"crypto/tls"
	"net"
)


type Client struct {
	conn interface {
		Read(b []byte) (n int, err error)
		Write(b []byte) (n int, err error)
		Close() error
		LocalAddr() net.Addr
		RemoteAddr() net.Addr
	}
	dataReader *DataReader
	dataWriter *DataWriter
}


// Dial makes a (un-secure) TELNET client connection to the system's 'loopback address'
// (also known as "localhost" or 127.0.0.1).
//
// If a secure connection is desired, use `DialTLS` instead.
func Dial() (*Client, error) {
	return DialTo("")
}

// DialTo makes a (un-secure) TELNET client connection to the the address specified by
// 'addr'.
//
// If a secure connection is desired, use `DialToTLS` instead.
func DialTo(addr string) (*Client, error) {

	const network = "tcp"

	if "" == addr {
		addr = "127.0.0.1:telnet"
	}

	conn, err := net.Dial(network, addr)
	if nil != err {
		return nil, err
	}

	dataReader := NewDataReader(conn)
	dataWriter := NewDataWriter(conn)

	client := Client{
		conn:conn,
		dataReader:dataReader,
		dataWriter:dataWriter,
	}

	return &client, nil
}


// DialTLS makes a (secure) TELNETS client connection to the system's 'loopback address'
// (also known as "localhost" or 127.0.0.1).
func DialTLS(tlsConfig *tls.Config) (*Client, error) {
	return DialToTLS("", tlsConfig)
}

// DialToTLS makes a (secure) TELNETS client connection to the the address specified by
// 'addr'.
func DialToTLS(addr string, tlsConfig *tls.Config) (*Client, error) {

	const network = "tcp"

	if "" == addr {
		addr = "127.0.0.1:telnets"
	}

	conn, err := tls.Dial(network, addr, tlsConfig)
	if nil != err {
		return nil, err
	}

	dataReader := NewDataReader(conn)
	dataWriter := NewDataWriter(conn)

	client := Client{
		conn:conn,
		dataReader:dataReader,
		dataWriter:dataWriter,
	}

	return &client, nil
}



// Close closes the client connection.
//
// Typical usage might look like:
//
//	telnetsClient, err = telnet.DialToTLS(addr, tlsConfig)
//	if nil != err {
//		//@TODO: Handle error.
//		return err
//	}
//	defer telnetsClient.Close()
func (client *Client) Close() error {
	return client.conn.Close()
}


// Read receives `n` bytes sent from the server to the client,
// and "returns" into `p`.
//
// Note that Read can only be used for receiving TELNET (and TELNETS) data from the server.
//
// TELNET (and TELNETS) command codes cannot be received using this method, as Read deals
// with TELNET (and TELNETS) "unescaping", and (when appropriate) filters out TELNET (and TELNETS)
// command codes.
//
// Read makes Client fit the io.Reader interface.
func (client *Client) Read(p []byte) (n int, err error) {
	return client.dataReader.Read(p)
}


// Write sends `n` bytes from 'p' to the server.
//
// Note that Write can only be used for sending TELNET (and TELNETS) data to the server.
//
// TELNET (and TELNETS) command codes cannot be sent using this method, as Write deals with
// TELNET (and TELNETS) "escaping", and will properly "escape" anything written with it.
//
// Write makes Client fit the io.Writer interface.
func (client *Client) Write(p []byte) (n int, err error) {
	return client.dataWriter.Write(p)
}


// LocalAddr returns the local network address.
func (client *Client) LocalAddr() net.Addr {
	return client.conn.LocalAddr()
}


// RemoteAddr returns the remote network address.
func (client *Client) RemoteAddr() net.Addr {
	return client.conn.RemoteAddr()
}
