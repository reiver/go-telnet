package telnet

import (
	"github.com/reiver/go-oi"

	"bytes"
	"io/ioutil"

	"testing"
)

func TestStandardCallerFromClientToServer(t *testing.T) {

	tests := []struct {
		Bytes    []byte
		Expected []byte
	}{
		{
			Bytes:    []byte{},
			Expected: []byte{},
		},

		{
			Bytes:    []byte("a"),
			Expected: []byte(""),
		},
		{
			Bytes:    []byte("b"),
			Expected: []byte(""),
		},
		{
			Bytes:    []byte("c"),
			Expected: []byte(""),
		},

		{
			Bytes:    []byte("a\n"),
			Expected: []byte("a\r\n"),
		},
		{
			Bytes:    []byte("b\n"),
			Expected: []byte("b\r\n"),
		},
		{
			Bytes:    []byte("c\n"),
			Expected: []byte("c\r\n"),
		},

		{
			Bytes:    []byte("a\nb\nc"),
			Expected: []byte("a\r\nb\r\n"),
		},

		{
			Bytes:    []byte("a\nb\nc\n"),
			Expected: []byte("a\r\nb\r\nc\r\n"),
		},

		{
			Bytes:    []byte("apple"),
			Expected: []byte(""),
		},
		{
			Bytes:    []byte("banana"),
			Expected: []byte(""),
		},
		{
			Bytes:    []byte("cherry"),
			Expected: []byte(""),
		},

		{
			Bytes:    []byte("apple\n"),
			Expected: []byte("apple\r\n"),
		},
		{
			Bytes:    []byte("banana\n"),
			Expected: []byte("banana\r\n"),
		},
		{
			Bytes:    []byte("cherry\n"),
			Expected: []byte("cherry\r\n"),
		},

		{
			Bytes:    []byte("apple\nbanana\ncherry"),
			Expected: []byte("apple\r\nbanana\r\n"),
		},

		{
			Bytes:    []byte("apple\nbanana\ncherry\n"),
			Expected: []byte("apple\r\nbanana\r\ncherry\r\n"),
		},

		{
			Bytes:    []byte("apple banana cherry"),
			Expected: []byte(""),
		},

		{
			Bytes:    []byte("apple banana cherry\n"),
			Expected: []byte("apple banana cherry\r\n"),
		},

		{
			Bytes:    []byte{255},
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 255},
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 255, 255},
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 255, 255, 255},
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 255, 255, 255, 255},
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 255, 255, 255, 255, 255},
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 255, 255, 255, 255, 255, 255},
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 255, 255, 255, 255, 255, 255, 255},
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 255, 255, 255, 255, 255, 255, 255, 255},
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			Expected: []byte{},
		},

		{
			Bytes:    []byte{255, '\n'},
			Expected: []byte{255, 255, '\r', '\n'},
		},
		{
			Bytes:    []byte{255, 255, '\n'},
			Expected: []byte{255, 255, 255, 255, '\r', '\n'},
		},
		{
			Bytes:    []byte{255, 255, 255, '\n'},
			Expected: []byte{255, 255, 255, 255, 255, 255, '\r', '\n'},
		},
		{
			Bytes:    []byte{255, 255, 255, 255, '\n'},
			Expected: []byte{255, 255, 255, 255, 255, 255, 255, 255, '\r', '\n'},
		},
		{
			Bytes:    []byte{255, 255, 255, 255, 255, '\n'},
			Expected: []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, '\r', '\n'},
		},
		{
			Bytes:    []byte{255, 255, 255, 255, 255, 255, '\n'},
			Expected: []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, '\r', '\n'},
		},
		{
			Bytes:    []byte{255, 255, 255, 255, 255, 255, 255, '\n'},
			Expected: []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, '\r', '\n'},
		},
		{
			Bytes:    []byte{255, 255, 255, 255, 255, 255, 255, 255, '\n'},
			Expected: []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, '\r', '\n'},
		},
		{
			Bytes:    []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, '\n'},
			Expected: []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, '\r', '\n'},
		},
		{
			Bytes:    []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, '\n'},
			Expected: []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, '\r', '\n'},
		},

		{
			Bytes:    []byte("apple\xff\xffbanana\xff\xffcherry"),
			Expected: []byte(""),
		},
		{
			Bytes:    []byte("\xff\xffapple\xff\xffbanana\xff\xffcherry\xff\xff"),
			Expected: []byte(""),
		},

		{
			Bytes:    []byte("apple\xffbanana\xffcherry\n"),
			Expected: []byte("apple\xff\xffbanana\xff\xffcherry\r\n"),
		},
		{
			Bytes:    []byte("\xffapple\xffbanana\xffcherry\xff\n"),
			Expected: []byte("\xff\xffapple\xff\xffbanana\xff\xffcherry\xff\xff\r\n"),
		},

		{
			Bytes:    []byte("apple\xff\xff\xff\xffbanana\xff\xff\xff\xffcherry"),
			Expected: []byte(""),
		},
		{
			Bytes:    []byte("\xff\xff\xff\xffapple\xff\xff\xff\xffbanana\xff\xff\xff\xffcherry\xff\xff\xff\xff"),
			Expected: []byte(""),
		},

		{
			Bytes:    []byte("apple\xff\xffbanana\xff\xffcherry\n"),
			Expected: []byte("apple\xff\xff\xff\xffbanana\xff\xff\xff\xffcherry\r\n"),
		},
		{
			Bytes:    []byte("\xff\xffapple\xff\xffbanana\xff\xffcherry\xff\xff\n"),
			Expected: []byte("\xff\xff\xff\xffapple\xff\xff\xff\xffbanana\xff\xff\xff\xffcherry\xff\xff\xff\xff\r\n"),
		},

		{
			Bytes:    []byte{255, 251, 24}, // IAC WILL TERMINAL-TYPE
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 252, 24}, // IAC WON'T TERMINAL-TYPE
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 253, 24}, // IAC DO TERMINAL-TYPE
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 254, 24}, // IAC DON'T TERMINAL-TYPE
			Expected: []byte{},
		},

		{
			Bytes:    []byte{255, 251, 24, '\n'}, // IAC WILL TERMINAL-TYPE '\n'
			Expected: []byte{255, 255, 251, 24, '\r', '\n'},
		},
		{
			Bytes:    []byte{255, 252, 24, '\n'}, // IAC WON'T TERMINAL-TYPE '\n'
			Expected: []byte{255, 255, 252, 24, '\r', '\n'},
		},
		{
			Bytes:    []byte{255, 253, 24, '\n'}, // IAC DO TERMINAL-TYPE '\n'
			Expected: []byte{255, 255, 253, 24, '\r', '\n'},
		},
		{
			Bytes:    []byte{255, 254, 24, '\n'}, // IAC DON'T TERMINAL-TYPE '\n'
			Expected: []byte{255, 255, 254, 24, '\r', '\n'},
		},

		{
			Bytes:    []byte{67, 255, 251, 24}, // 'C' IAC WILL TERMINAL-TYPE
			Expected: []byte{},
		},
		{
			Bytes:    []byte{67, 255, 252, 24}, // 'C' IAC WON'T TERMINAL-TYPE
			Expected: []byte{},
		},
		{
			Bytes:    []byte{67, 255, 253, 24}, // 'C' IAC DO TERMINAL-TYPE
			Expected: []byte{},
		},
		{
			Bytes:    []byte{67, 255, 254, 24}, // 'C' IAC DON'T TERMINAL-TYPE
			Expected: []byte{},
		},

		{
			Bytes:    []byte{67, 255, 251, 24, '\n'}, // 'C' IAC WILL TERMINAL-TYPE '\n'
			Expected: []byte{67, 255, 255, 251, 24, '\r', '\n'},
		},
		{
			Bytes:    []byte{67, 255, 252, 24, '\n'}, // 'C' IAC WON'T TERMINAL-TYPE '\n'
			Expected: []byte{67, 255, 255, 252, 24, '\r', '\n'},
		},
		{
			Bytes:    []byte{67, 255, 253, 24, '\n'}, // 'C' IAC DO TERMINAL-TYPE '\n'
			Expected: []byte{67, 255, 255, 253, 24, '\r', '\n'},
		},
		{
			Bytes:    []byte{67, 255, 254, 24, '\n'}, // 'C' IAC DON'T TERMINAL-TYPE '\n'
			Expected: []byte{67, 255, 255, 254, 24, '\r', '\n'},
		},

		{
			Bytes:    []byte{255, 251, 24, 68}, // IAC WILL TERMINAL-TYPE 'D'
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 252, 24, 68}, // IAC WON'T TERMINAL-TYPE 'D'
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 253, 24, 68}, // IAC DO TERMINAL-TYPE 'D'
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 254, 24, 68}, // IAC DON'T TERMINAL-TYPE 'D'
			Expected: []byte{},
		},

		{
			Bytes:    []byte{255, 251, 24, 68, '\n'}, // IAC WILL TERMINAL-TYPE 'D' '\n'
			Expected: []byte{255, 255, 251, 24, 68, '\r', '\n'},
		},
		{
			Bytes:    []byte{255, 252, 24, 68, '\n'}, // IAC WON'T TERMINAL-TYPE 'D' '\n'
			Expected: []byte{255, 255, 252, 24, 68, '\r', '\n'},
		},
		{
			Bytes:    []byte{255, 253, 24, 68, '\n'}, // IAC DO TERMINAL-TYPE 'D' '\n'
			Expected: []byte{255, 255, 253, 24, 68, '\r', '\n'},
		},
		{
			Bytes:    []byte{255, 254, 24, 68, '\n'}, // IAC DON'T TERMINAL-TYPE 'D' '\n'
			Expected: []byte{255, 255, 254, 24, 68, '\r', '\n'},
		},

		{
			Bytes:    []byte{67, 255, 251, 24, 68}, // 'C' IAC WILL TERMINAL-TYPE 'D'
			Expected: []byte{},
		},
		{
			Bytes:    []byte{67, 255, 252, 24, 68}, // 'C' IAC WON'T TERMINAL-TYPE 'D'
			Expected: []byte{},
		},
		{
			Bytes:    []byte{67, 255, 253, 24, 68}, // 'C' IAC DO TERMINAL-TYPE 'D'
			Expected: []byte{},
		},
		{
			Bytes:    []byte{67, 255, 254, 24, 68}, // 'C' IAC DON'T TERMINAL-TYPE 'D'
			Expected: []byte{},
		},

		{
			Bytes:    []byte{67, 255, 251, 24, 68, '\n'}, // 'C' IAC WILL TERMINAL-TYPE 'D' '\n'
			Expected: []byte{67, 255, 255, 251, 24, 68, '\r', '\n'},
		},
		{
			Bytes:    []byte{67, 255, 252, 24, 68, '\n'}, // 'C' IAC WON'T TERMINAL-TYPE 'D' '\n'
			Expected: []byte{67, 255, 255, 252, 24, 68, '\r', '\n'},
		},
		{
			Bytes:    []byte{67, 255, 253, 24, 68, '\n'}, // 'C' IAC DO TERMINAL-TYPE 'D' '\n'
			Expected: []byte{67, 255, 255, 253, 24, 68, '\r', '\n'},
		},
		{
			Bytes:    []byte{67, 255, 254, 24, 68, '\n'}, // 'C' IAC DON'T TERMINAL-TYPE 'D' '\n'
			Expected: []byte{67, 255, 255, 254, 24, 68, '\r', '\n'},
		},

		{
			Bytes:    []byte{255, 250, 24, 1, 255, 240}, // IAC SB TERMINAL-TYPE SEND IAC SE
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 250, 24, 0, 68, 69, 67, 45, 86, 84, 53, 50, 255, 240}, // IAC SB TERMINAL-TYPE IS "DEC-VT52" IAC SE
			Expected: []byte{},
		},

		{
			Bytes:    []byte{255, 250, 24, 1, 255, 240, '\n'}, // IAC SB TERMINAL-TYPE SEND IAC SE '\n'
			Expected: []byte{255, 255, 250, 24, 1, 255, 255, 240, '\r', '\n'},
		},
		{
			Bytes:    []byte{255, 250, 24, 0, 68, 69, 67, 45, 86, 84, 53, 50, 255, 240, '\n'}, // IAC SB TERMINAL-TYPE IS "DEC-VT52" IAC SE '\n'
			Expected: []byte{255, 255, 250, 24, 0, 68, 69, 67, 45, 86, 84, 53, 50, 255, 255, 240, '\r', '\n'},
		},

		{
			Bytes:    []byte{67, 255, 250, 24, 1, 255, 240}, // 'C' IAC SB TERMINAL-TYPE SEND IAC SE
			Expected: []byte{},
		},
		{
			Bytes:    []byte{67, 255, 250, 24, 0, 68, 69, 67, 45, 86, 84, 53, 50, 255, 240}, // 'C' IAC SB TERMINAL-TYPE IS "DEC-VT52" IAC SE
			Expected: []byte{},
		},

		{
			Bytes:    []byte{67, 255, 250, 24, 1, 255, 240, '\n'}, // 'C' IAC SB TERMINAL-TYPE SEND IAC SE '\n'
			Expected: []byte{67, 255, 255, 250, 24, 1, 255, 255, 240, '\r', '\n'},
		},
		{
			Bytes:    []byte{67, 255, 250, 24, 0, 68, 69, 67, 45, 86, 84, 53, 50, 255, 240, '\n'}, // 'C' IAC SB TERMINAL-TYPE IS "DEC-VT52" IAC SE '\n'
			Expected: []byte{67, 255, 255, 250, 24, 0, 68, 69, 67, 45, 86, 84, 53, 50, 255, 255, 240, '\r', '\n'},
		},

		{
			Bytes:    []byte{255, 250, 24, 1, 255, 240, 68}, // IAC SB TERMINAL-TYPE SEND IAC SE 'D'
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 250, 24, 0, 68, 69, 67, 45, 86, 84, 53, 50, 255, 240, 68}, // IAC SB TERMINAL-TYPE IS "DEC-VT52" IAC SE 'D'
			Expected: []byte{},
		},

		{
			Bytes:    []byte{255, 250, 24, 1, 255, 240, 68, '\n'}, // IAC SB TERMINAL-TYPE SEND IAC SE 'D' '\n'
			Expected: []byte{255, 255, 250, 24, 1, 255, 255, 240, 68, '\r', '\n'},
		},
		{
			Bytes:    []byte{255, 250, 24, 0, 68, 69, 67, 45, 86, 84, 53, 50, 255, 240, 68, '\n'}, // IAC SB TERMINAL-TYPE IS "DEC-VT52" IAC SE 'D' '\n'
			Expected: []byte{255, 255, 250, 24, 0, 68, 69, 67, 45, 86, 84, 53, 50, 255, 255, 240, 68, '\r', '\n'},
		},

		{
			Bytes:    []byte{67, 255, 250, 24, 1, 255, 240, 68}, // 'C' IAC SB TERMINAL-TYPE SEND IAC SE 'D'
			Expected: []byte{},
		},
		{
			Bytes:    []byte{67, 255, 250, 24, 0, 68, 69, 67, 45, 86, 84, 53, 50, 255, 240, 68}, // 'C' IAC SB TERMINAL-TYPE IS "DEC-VT52" IAC SE 'D'
			Expected: []byte{},
		},

		{
			Bytes:    []byte{67, 255, 250, 24, 1, 255, 240, 68, '\n'}, // 'C' IAC SB TERMINAL-TYPE SEND IAC SE 'D' '\n'
			Expected: []byte{67, 255, 255, 250, 24, 1, 255, 255, 240, 68, '\r', '\n'},
		},
		{
			Bytes:    []byte{67, 255, 250, 24, 0, 68, 69, 67, 45, 86, 84, 53, 50, 255, 240, 68, '\n'}, // 'C' IAC SB TERMINAL-TYPE IS "DEC-VT52" IAC SE 'D' '\n'
			Expected: []byte{67, 255, 255, 250, 24, 0, 68, 69, 67, 45, 86, 84, 53, 50, 255, 255, 240, 68, '\r', '\n'},
		},

		{
			Bytes:    []byte{255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 255, 240}, // IAC SB 0 1 2 3 4 5 6 7 8 9 10 11 12 13 IAC SE
			Expected: []byte{255, 255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, '\r', 10},
		},
		{
			Bytes:    []byte{67, 255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 255, 240}, // 'C' IAC SB 0 1 2 3 4 5 6 7 8 9 10 11 12 13 IAC SE
			Expected: []byte{67, 255, 255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, '\r', 10},
		},
		{
			Bytes:    []byte{255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 255, 240, 68}, // IAC SB 0 1 2 3 4 5 6 7 8 9 10 11 12 13 IAC SE 'D'
			Expected: []byte{255, 255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, '\r', 10},
		},
		{
			Bytes:    []byte{67, 255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 255, 240, 68}, // 'C' IAC SB 0 1 2 3 4 5 6 7 8 9 10 11 12 13 IAC SE 'D'
			Expected: []byte{67, 255, 255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, '\r', 10},
		},

		{
			Bytes:    []byte{255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 255, 240, '\n'}, // IAC SB 0 1 2 3 4 5 6 7 8 9 10 11 12 13 IAC SE '\n'
			Expected: []byte{255, 255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, '\r', 10, 11, 12, 13, 255, 255, 240, '\r', '\n'},
		},
		{
			Bytes:    []byte{67, 255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 255, 240, '\n'}, // 'C' IAC SB 0 1 2 3 4 5 6 7 8 9 10 11 12 13 IAC SE '\n'
			Expected: []byte{67, 255, 255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, '\r', 10, 11, 12, 13, 255, 255, 240, '\r', '\n'},
		},
		{
			Bytes:    []byte{255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 255, 240, 68, '\n'}, // IAC SB 0 1 2 3 4 5 6 7 8 9 10 11 12 13 IAC SE 'D' '\n'
			Expected: []byte{255, 255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, '\r', 10, 11, 12, 13, 255, 255, 240, 68, '\r', '\n'},
		},
		{
			Bytes:    []byte{67, 255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 255, 240, 68, '\n'}, // 'C' IAC SB 0 1 2 3 4 5 6 7 8 9 10 11 12 13 IAC SE 'D' '\n'
			Expected: []byte{67, 255, 255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, '\r', 10, 11, 12, 13, 255, 255, 240, 68, '\r', '\n'},
		},
	}

	for testNumber, test := range tests {
		var stdinBuffer bytes.Buffer
		var stdoutBuffer bytes.Buffer
		var stderrBuffer bytes.Buffer

		stdinBuffer.Write(test.Bytes) // <----------------- The important difference between the 2 loops.

		stdin := ioutil.NopCloser(&stdinBuffer)
		stdout := oi.WriteNopCloser(&stdoutBuffer)
		stderr := oi.WriteNopCloser(&stderrBuffer)

		var ctx Context = nil

		var dataWriterBuffer bytes.Buffer
		dataWriter := newDataWriter(&dataWriterBuffer)

		dataReader := newDataReader(bytes.NewReader([]byte{})) // <----------------- The important difference between the 2 loops.

		standardCallerCallTELNET(stdin, stdout, stderr, ctx, dataWriter, dataReader)

		if expected, actual := string(test.Expected), dataWriterBuffer.String(); expected != actual {
			t.Errorf("For test #%d, expected %q, but actually got %q; for %q.", testNumber, expected, actual, test.Bytes)
			continue
		}

		if expected, actual := "", stdoutBuffer.String(); expected != actual {
			t.Errorf("For test #%d, expected %q, but actually got %q.", testNumber, expected, actual)
			continue
		}

		if expected, actual := "", stderrBuffer.String(); expected != actual {
			t.Errorf("For test #%d, expected %q, but actually got %q.", testNumber, expected, actual)
			continue
		}
	}
}

func TestStandardCallerFromServerToClient(t *testing.T) {

	tests := []struct {
		Bytes    []byte
		Expected []byte
	}{
		{
			Bytes:    []byte{},
			Expected: []byte{},
		},

		{
			Bytes:    []byte("a"),
			Expected: []byte("a"),
		},
		{
			Bytes:    []byte("b"),
			Expected: []byte("b"),
		},
		{
			Bytes:    []byte("c"),
			Expected: []byte("c"),
		},

		{
			Bytes:    []byte("apple"),
			Expected: []byte("apple"),
		},
		{
			Bytes:    []byte("banana"),
			Expected: []byte("banana"),
		},
		{
			Bytes:    []byte("cherry"),
			Expected: []byte("cherry"),
		},

		{
			Bytes:    []byte("apple banana cherry"),
			Expected: []byte("apple banana cherry"),
		},

		{
			Bytes:    []byte{255, 255},
			Expected: []byte{255},
		},
		{
			Bytes:    []byte{255, 255, 255, 255},
			Expected: []byte{255, 255},
		},
		{
			Bytes:    []byte{255, 255, 255, 255, 255, 255},
			Expected: []byte{255, 255, 255},
		},
		{
			Bytes:    []byte{255, 255, 255, 255, 255, 255, 255, 255},
			Expected: []byte{255, 255, 255, 255},
		},
		{
			Bytes:    []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			Expected: []byte{255, 255, 255, 255, 255},
		},

		{
			Bytes:    []byte("apple\xff\xffbanana\xff\xffcherry"),
			Expected: []byte("apple\xffbanana\xffcherry"),
		},
		{
			Bytes:    []byte("\xff\xffapple\xff\xffbanana\xff\xffcherry\xff\xff"),
			Expected: []byte("\xffapple\xffbanana\xffcherry\xff"),
		},

		{
			Bytes:    []byte("apple\xff\xff\xff\xffbanana\xff\xff\xff\xffcherry"),
			Expected: []byte("apple\xff\xffbanana\xff\xffcherry"),
		},
		{
			Bytes:    []byte("\xff\xff\xff\xffapple\xff\xff\xff\xffbanana\xff\xff\xff\xffcherry\xff\xff\xff\xff"),
			Expected: []byte("\xff\xffapple\xff\xffbanana\xff\xffcherry\xff\xff"),
		},

		{
			Bytes:    []byte{255, 251, 24}, // IAC WILL TERMINAL-TYPE
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 252, 24}, // IAC WON'T TERMINAL-TYPE
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 253, 24}, // IAC DO TERMINAL-TYPE
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 254, 24}, // IAC DON'T TERMINAL-TYPE
			Expected: []byte{},
		},

		{
			Bytes:    []byte{67, 255, 251, 24}, // 'C' IAC WILL TERMINAL-TYPE
			Expected: []byte{67},
		},
		{
			Bytes:    []byte{67, 255, 252, 24}, // 'C' IAC WON'T TERMINAL-TYPE
			Expected: []byte{67},
		},
		{
			Bytes:    []byte{67, 255, 253, 24}, // 'C' IAC DO TERMINAL-TYPE
			Expected: []byte{67},
		},
		{
			Bytes:    []byte{67, 255, 254, 24}, // 'C' IAC DON'T TERMINAL-TYPE
			Expected: []byte{67},
		},

		{
			Bytes:    []byte{255, 251, 24, 68}, // IAC WILL TERMINAL-TYPE 'D'
			Expected: []byte{68},
		},
		{
			Bytes:    []byte{255, 252, 24, 68}, // IAC WON'T TERMINAL-TYPE 'D'
			Expected: []byte{68},
		},
		{
			Bytes:    []byte{255, 253, 24, 68}, // IAC DO TERMINAL-TYPE 'D'
			Expected: []byte{68},
		},
		{
			Bytes:    []byte{255, 254, 24, 68}, // IAC DON'T TERMINAL-TYPE 'D'
			Expected: []byte{68},
		},

		{
			Bytes:    []byte{67, 255, 251, 24, 68}, // 'C' IAC WILL TERMINAL-TYPE 'D'
			Expected: []byte{67, 68},
		},
		{
			Bytes:    []byte{67, 255, 252, 24, 68}, // 'C' IAC WON'T TERMINAL-TYPE 'D'
			Expected: []byte{67, 68},
		},
		{
			Bytes:    []byte{67, 255, 253, 24, 68}, // 'C' IAC DO TERMINAL-TYPE 'D'
			Expected: []byte{67, 68},
		},
		{
			Bytes:    []byte{67, 255, 254, 24, 68}, // 'C' IAC DON'T TERMINAL-TYPE 'D'
			Expected: []byte{67, 68},
		},

		{
			Bytes:    []byte{255, 250, 24, 1, 255, 240}, // IAC SB TERMINAL-TYPE SEND IAC SE
			Expected: []byte{},
		},
		{
			Bytes:    []byte{255, 250, 24, 0, 68, 69, 67, 45, 86, 84, 53, 50, 255, 240}, // IAC SB TERMINAL-TYPE IS "DEC-VT52" IAC SE
			Expected: []byte{},
		},

		{
			Bytes:    []byte{67, 255, 250, 24, 1, 255, 240}, // 'C' IAC SB TERMINAL-TYPE SEND IAC SE
			Expected: []byte{67},
		},
		{
			Bytes:    []byte{67, 255, 250, 24, 0, 68, 69, 67, 45, 86, 84, 53, 50, 255, 240}, // 'C' IAC SB TERMINAL-TYPE IS "DEC-VT52" IAC SE
			Expected: []byte{67},
		},

		{
			Bytes:    []byte{255, 250, 24, 1, 255, 240, 68}, // IAC SB TERMINAL-TYPE SEND IAC SE 'D'
			Expected: []byte{68},
		},
		{
			Bytes:    []byte{255, 250, 24, 0, 68, 69, 67, 45, 86, 84, 53, 50, 255, 240, 68}, // IAC SB TERMINAL-TYPE IS "DEC-VT52" IAC SE 'D'
			Expected: []byte{68},
		},

		{
			Bytes:    []byte{67, 255, 250, 24, 1, 255, 240, 68}, // 'C' IAC SB TERMINAL-TYPE SEND IAC SE 'D'
			Expected: []byte{67, 68},
		},
		{
			Bytes:    []byte{67, 255, 250, 24, 0, 68, 69, 67, 45, 86, 84, 53, 50, 255, 240, 68}, // 'C' IAC SB TERMINAL-TYPE IS "DEC-VT52" IAC SE 'D'
			Expected: []byte{67, 68},
		},

		{
			Bytes:    []byte{255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 255, 240}, // IAC SB 0 1 2 3 4 5 6 7 8 9 10 11 12 13 IAC SE
			Expected: []byte{},
		},
		{
			Bytes:    []byte{67, 255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 255, 240}, // 'C' IAC SB 0 1 2 3 4 5 6 7 8 9 10 11 12 13 IAC SE
			Expected: []byte{67},
		},
		{
			Bytes:    []byte{255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 255, 240, 68}, // IAC SB 0 1 2 3 4 5 6 7 8 9 10 11 12 13 IAC SE 'D'
			Expected: []byte{68},
		},
		{
			Bytes:    []byte{67, 255, 250, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 255, 240, 68}, // 'C' IAC SB 0 1 2 3 4 5 6 7 8 9 10 11 12 13 IAC SE 'D'
			Expected: []byte{67, 68},
		},

		//@TODO: Is this correct? Can IAC appear between thee 'IAC SB' and ''IAC SE'?... and if "yes", do escaping rules apply?
		{
			Bytes:    []byte{255, 250, 255, 255, 240, 255, 240}, //     IAC SB 255 255 240 IAC SE = IAC SB IAC IAC SE IAC SE
			Expected: []byte{},
		},
		{
			Bytes:    []byte{67, 255, 250, 255, 255, 240, 255, 240}, // 'C' IAC SB 255 255 240 IAC SE = IAC SB IAC IAC SE IAC SE
			Expected: []byte{67},
		},
		{
			Bytes:    []byte{255, 250, 255, 255, 240, 255, 240, 68}, //     IAC SB 255 255 240 IAC SE = IAC SB IAC IAC SE IAC SE 'D'
			Expected: []byte{68},
		},
		{
			Bytes:    []byte{67, 255, 250, 255, 255, 240, 255, 240, 68}, // 'C' IAC SB 255 255 240 IAC SE = IAC SB IAC IAC SE IAC SE 'D'
			Expected: []byte{67, 68},
		},

		//@TODO: Is this correct? Can IAC appear between thee 'IAC SB' and ''IAC SE'?... and if "yes", do escaping rules apply?
		{
			Bytes:    []byte{255, 250, 71, 255, 255, 240, 255, 240}, //     IAC SB 'G' 255 255 240 IAC SE = IAC SB 'G' IAC IAC SE IAC SE
			Expected: []byte{},
		},
		{
			Bytes:    []byte{67, 255, 250, 71, 255, 255, 240, 255, 240}, // 'C' IAC SB 'G' 255 255 240 IAC SE = IAC SB 'G' IAC IAC SE IAC SE
			Expected: []byte{67},
		},
		{
			Bytes:    []byte{255, 250, 71, 255, 255, 240, 255, 240, 68}, //     IAC SB 'G' 255 255 240 IAC SE = IAC SB 'G' IAC IAC SE IAC SE 'D'
			Expected: []byte{68},
		},
		{
			Bytes:    []byte{67, 255, 250, 71, 255, 255, 240, 255, 240, 68}, // 'C' IAC SB 'G' 255 255 240 IAC SE = IAC 'G' SB IAC IAC SE IAC SE 'D'
			Expected: []byte{67, 68},
		},

		//@TODO: Is this correct? Can IAC appear between thee 'IAC SB' and ''IAC SE'?... and if "yes", do escaping rules apply?
		{
			Bytes:    []byte{255, 250, 255, 255, 240, 72, 255, 240}, //     IAC SB 255 255 240 'H' IAC SE = IAC SB IAC IAC SE 'H' IAC SE
			Expected: []byte{},
		},
		{
			Bytes:    []byte{67, 255, 250, 255, 255, 240, 72, 255, 240}, // 'C' IAC SB 255 255 240 'H' IAC SE = IAC SB IAC IAC SE 'H' IAC SE
			Expected: []byte{67},
		},
		{
			Bytes:    []byte{255, 250, 255, 255, 240, 72, 255, 240, 68}, //     IAC SB 255 255 240 'H' IAC SE = IAC SB IAC IAC SE 'H' IAC SE 'D'
			Expected: []byte{68},
		},
		{
			Bytes:    []byte{67, 255, 250, 255, 255, 240, 72, 255, 240, 68}, // 'C' IAC SB 255 255 240 'H' IAC SE = IAC SB IAC IAC SE 'H' IAC SE 'D'
			Expected: []byte{67, 68},
		},

		//@TODO: Is this correct? Can IAC appear between thee 'IAC SB' and ''IAC SE'?... and if "yes", do escaping rules apply?
		{
			Bytes:    []byte{255, 250, 71, 255, 255, 240, 72, 255, 240}, //     IAC SB 'G' 255 255 240 'H' IAC SE = IAC SB 'G' IAC IAC SE 'H' IAC SE
			Expected: []byte{},
		},
		{
			Bytes:    []byte{67, 255, 250, 71, 255, 255, 240, 72, 255, 240}, // 'C' IAC SB 'G' 255 255 240 'H' IAC SE = IAC SB 'G' IAC IAC SE 'H' IAC SE
			Expected: []byte{67},
		},
		{
			Bytes:    []byte{255, 250, 71, 255, 255, 240, 72, 255, 240, 68}, //     IAC SB 'G' 255 255 240 'H' IAC SE = IAC SB 'G' IAC IAC SE 'H' IAC SE 'D'
			Expected: []byte{68},
		},
		{
			Bytes:    []byte{67, 255, 250, 71, 255, 255, 240, 72, 255, 240, 68}, // 'C' IAC SB 'G' 255 255 240 'H' IAC SE = IAC 'G' SB IAC IAC SE 'H' IAC SE 'D'
			Expected: []byte{67, 68},
		},
	}

	for testNumber, test := range tests {
		var stdinBuffer bytes.Buffer
		var stdoutBuffer bytes.Buffer
		var stderrBuffer bytes.Buffer

		stdin := ioutil.NopCloser(&stdinBuffer)
		stdout := oi.WriteNopCloser(&stdoutBuffer)
		stderr := oi.WriteNopCloser(&stderrBuffer)

		var ctx Context = nil

		var dataWriterBuffer bytes.Buffer
		dataWriter := newDataWriter(&dataWriterBuffer)

		dataReader := newDataReader(bytes.NewReader(test.Bytes)) // <----------------- The important difference between the 2 loops.

		standardCallerCallTELNET(stdin, stdout, stderr, ctx, dataWriter, dataReader)

		if expected, actual := "", dataWriterBuffer.String(); expected != actual {
			t.Errorf("For test #%d, expected %q, but actually got %q; for %q.", testNumber, expected, actual, test.Bytes)
			continue
		}

		if expected, actual := string(test.Expected), stdoutBuffer.String(); expected != actual {
			t.Errorf("For test #%d, expected %q, but actually got %q.", testNumber, expected, actual)
			continue
		}

		if expected, actual := "", stderrBuffer.String(); expected != actual {
			t.Errorf("For test #%d, expected %q, but actually got %q.", testNumber, expected, actual)
			continue
		}

	}
}
