package telnet


// A Caller represents the client end of a TELNET (or TELNETS) connection.
//
// Writing data to the *DataWriter passed as an argument to the CallTELNET method
// will send data to the TELNET (or TELNETS) server.
//
// Reading data from the *DataReadere passed as an argument to the CallTELNET method
// will receive data from the TELNET server.
//
// The *DataWriter's Write method sends "escaped" TELNET (and TELNETS) data.
//
// The *DataReader's Reader method receives "un-escaped" TELNET (and TELNETS) data.
//
// Meaning, for example, that it does NOT include TELNET (and TELNETS) command sequences.
type Caller interface {
	CallTELNET(Context, *DataWriter, *DataReader)
}
