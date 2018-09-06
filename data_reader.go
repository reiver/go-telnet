package telnet

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

var (
	errCorrupted = errors.New("Corrupted")
)

const (
	IAC = 255

	SB = 250
	SE = 240

	WILL = 251
	WONT = 252
	DO   = 253
	DONT = 254
)

// An internalDataReader deals with "un-escaping" according to the TELNET protocol.
//
// In the TELNET protocol byte value 255 is special.
//
// The TELNET protocol calls byte value 255: "IAC". Which is short for "interpret as command".
//
// The TELNET protocol also has a distinction between 'data' and 'commands'.
//
//(DataReader is targetted toward TELNET 'data', not TELNET 'commands'.)
//
// If a byte with value 255 (=IAC) appears in the data, then it must be escaped.
//
// Escaping byte value 255 (=IAC) in the data is done by putting 2 of them in a row.
//
// So, for example:
//
//	[]byte{255} -> []byte{255, 255}
//
// Or, for a more complete example, if we started with the following:
//
//	[]byte{1, 55, 2, 155, 3, 255, 4, 40, 255, 30, 20}
//
// ... TELNET escaping would produce the following:
//
//	[]byte{1, 55, 2, 155, 3, 255, 255, 4, 40, 255, 255, 30, 20}
//
// (Notice that each "255" in the original byte array became 2 "255"s in a row.)
//
// DataReader deals with "un-escaping". In other words, it un-does what was shown
// in the examples.
//
// So, for example, it does this:
//
//	[]byte{255, 255} -> []byte{255}
//
// And, for example, goes from this:
//
//	[]byte{1, 55, 2, 155, 3, 255, 255, 4, 40, 255, 255, 30, 20}
//
// ... to this:
//
//	[]byte{1, 55, 2, 155, 3, 255, 4, 40, 255, 30, 20}
type internalDataReader struct {
	wrapped io.Reader
	state   state
}

// newDataReader creates a new DataReader reading from 'r'.
func newDataReader(r io.Reader) *internalDataReader {
	reader := internalDataReader{
		wrapped: r,
		state:   copyData,
	}

	return &reader
}

// Read reads the TELNET escaped data from the wrapped io.Reader, and "un-escapes" it into 'data'.
// It executes exactly one Read on the underlying reader every time it is called.
// Callers should be careful to truncate data to the number of bytes read,
// since this reader is expected to drop bytes from the underlying reader
// as described above when required to by the TELNET protocol.
func (r *internalDataReader) Read(data []byte) (int, error) {
	mach := &machine{
		from: make([]byte, len(data)),
		to:   data,
	}
	n, err := r.wrapped.Read(mach.from)
	if err != nil {
		return 0, err
	}
	mach.from = mach.from[:n]
	for err == nil && mach.InputRemaining() {
		r.state, err = r.state(mach)
	}
	return mach.written, err
}

// Unescaping of data read from the underlying reader is done using
// a state machine so it is resumable across reads.
type machine struct {
	from, to      []byte
	read, written int
}

// Index returns the offset from the read pointer of the first occurrence
// of the byte b, or -1 if that byte is not found in the remainder of the
// data read from the underlying reader.
func (m *machine) Index(b byte) int {
	return bytes.Index(m.from[m.read:], []byte{b})
}

// Copy copies up to n bytes from the underlying reader to the destination
// buffer, advancing both the read and write pointers by this amount.
func (m *machine) Copy(n int) {
	// Deliberately no bounds check here, because asking this
	// code to read past the end of m.from should never happen.
	copied := copy(m.to[m.written:], m.from[m.read:m.read+n])
	m.written += copied
	m.read += copied
}

// WriteByte writes the provided byte to the destnation buffer and advances
// the write pointer.
func (m *machine) WriteByte(b byte) {
	m.to[m.written] = b
	m.written++
}

// ConsumeByte reads and returns the next byte from the underlying reader,
// advancing the read pointer.
func (m *machine) ConsumeByte() byte {
	b := m.from[m.read]
	m.read++
	return b
}

// InputRemaining returns true as long as there is still data available
// to read from the input buffer.
func (m *machine) InputRemaining() bool {
	if m.read >= len(m.from) {
		return false
	}
	return true
}

// State machine states are functions that take the machine and
// return new states and optionally errors.
type state func(*machine) (state, error)

// The copyData state copies data from machine.from to machine.to
// until it encounters an IAC byte or the end of from.
func copyData(mach *machine) (state, error) {
	idx := mach.Index(IAC)
	if idx < 0 {
		// No escape bytes, so just copy remaining data and return.
		mach.Copy(len(mach.from) - mach.read)
		return copyData, nil
	}
	// Copy data up to IAC.
	mach.Copy(idx)
	return consumeIAC, nil
}

// The consumeIAC state eats an IAC byte and returns consumeCmd.
func consumeIAC(mach *machine) (state, error) {
	if b := mach.ConsumeByte(); b != IAC {
		return copyData, fmt.Errorf("expected IAC byte, got %c", b)
	}
	return consumeCmd, nil
}

// The consumeCmd state eats one of the known telnet command bytes.
func consumeCmd(mach *machine) (state, error) {
	switch b := mach.ConsumeByte(); b {
	case WILL, WONT, DO, DONT:
		// WILL, WONT, DO and DONT have an extra command byte
		// that shouldn't make it to the output slice.
		// We need to consume it before going back to copying data.
		return consumeWWDD, nil
	case IAC:
		// IAC IAC => un-escape; write IAC to output
		// and go back to copying data.
		mach.WriteByte(IAC)
		return copyData, nil
	case SB:
		// IAC SB => switch to consuming status.
		return consumeStatus, nil
	case SE:
		// IAC SE => go back to copying data.
		return copyData, nil
	default:
		// IAC <other bytes> is a protocol error.
		return copyData, fmt.Errorf("expected command byte, got %c", b)
	}
}

// The consumeWWDD state eats one byte then resumes copying data.
func consumeWWDD(mach *machine) (state, error) {
	mach.ConsumeByte()
	return copyData, nil
}

// The consumeStatus state eats data until it encounters an IAC
// byte or the end of from.
func consumeStatus(mach *machine) (state, error) {
	// We don't try to understand the status commands,
	// we just strip them from the output, which means
	// dropping input data until we read IAC SE.
	idx := mach.Index(IAC)
	if idx < 0 {
		// No escape bytes, so just skip remaining data and return.
		mach.read = len(mach.from)
		return consumeStatus, nil
	}
	// Skip up to IAC.
	mach.read += idx
	return consumeStatusIAC, nil
}

// The consumeStatusIAC state eats an IAC byte and returns consumeStatusCmd.
func consumeStatusIAC(mach *machine) (state, error) {
	if b := mach.ConsumeByte(); b != IAC {
		return consumeStatus, fmt.Errorf("expected IAC byte, got %c", b)
	}
	return consumeStatusCmd, nil
}

// The consumeStatusCmd state eats a byte. If that byte is SE the machine
// goes back to copying data, otherwise it goes back to consuming status.
func consumeStatusCmd(mach *machine) (state, error) {
	switch b := mach.ConsumeByte(); b {
	case SE:
		// IAC SE => go back to copying data normally
		return copyData, nil
	default:
		// IAC <byte> => continue eating SB
		return consumeStatus, nil
	}
}
