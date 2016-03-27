package telnet


import (
	"io"
)


// A Handler serves a TELNET (or TELNETS) connection.
//
// Writing data to the io.Writer passed as an argument the ServeTELNET method
// will send data to the TELNET (or TELNETS) client.
//
// Reading data from the io.Reader passed as an argument to the ServeTELNET method
// will receive data from the TELNET (or TELNETS) client.
//
// The Reader's Read method and Writer's Write method receive and send
// the raw TELNET (or TELNETS) communications!
//
// Meaning, for example, that it includes commands and escaping.
//
// You can use DataReader to un-escape raw TELNET (or TELNETS) communications.
//
// And you can use DataWriter to escape into raw TELNET (or TELNETS) communications.
type Handler interface {
	ServeTELNET(Context, io.Writer, io.Reader)
}
