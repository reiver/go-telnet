package telnet


import (
	"io"
)


// A Handler serves a TELNET connection.
//
// Writing data to the io.Writer passed as an argument the ServeTELNET method
// will send data to the TELNET client.
//
// Reading data from the io.Reader passed as an argument the ServeTELNET method
// will receive data from the TELNET client.
//
// The Reader's Read method and Writer's Write method receive and send
// the raw TELNET communications!
//
// Meaning, for example, that it includes commands and escaping.
//
// You could use DataReader to un-escape raw TELNET communications.
//
// And you could use DataWriter to escape into raw TELNET communications.
type Handler interface {
	ServeTELNET(Context, io.Writer, io.Reader)
}
