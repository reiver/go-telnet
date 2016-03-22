package telnet


import (
	"io"
)


// EchoHandler is a simple TELNET server which "echos" back to the client any (non-command)
// data back to the TELNET client, it received from the TELNET client.
var EchoHandler Handler = internalEchoHandler{}


type internalEchoHandler struct{}


func (handler internalEchoHandler) ServeTELNET(ctx Context, w io.Writer, r io.Reader) {

	writer := NewDataWriter(w)
	reader := NewDataReader(r)

	//@TODO: Will this overflow somehow if the 'written' goes past the int64 limit?
	_, err := io.Copy(writer, reader)
	if nil != err {
		//@TODO: Is panic()ing what we really want to do here?
		panic(err)
	}
}
