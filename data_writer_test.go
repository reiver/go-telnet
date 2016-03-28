package telnet


import (
	"bytes"

	"testing"
)


func TestDataWriter(t *testing.T) {

	tests := []struct{
		Bytes    []byte
		Expected []byte
	}{
		{
			Bytes:    []byte{},
			Expected: []byte{},
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
			Bytes:    []byte{255},
			Expected: []byte{255,255},
		},
		{
			Bytes:    []byte{255,255},
			Expected: []byte{255,255,255,255},
		},
		{
			Bytes:    []byte{255,255,255},
			Expected: []byte{255,255,255,255,255,255},
		},
		{
			Bytes:    []byte{255,255,255,255},
			Expected: []byte{255,255,255,255,255,255,255,255},
		},
		{
			Bytes:    []byte{255,255,255,255,255},
			Expected: []byte{255,255,255,255,255,255,255,255,255,255},
		},



		{
			Bytes:    []byte("apple\xffbanana\xffcherry"),
			Expected: []byte("apple\xff\xffbanana\xff\xffcherry"),
		},
		{
			Bytes:    []byte("\xffapple\xffbanana\xffcherry\xff"),
			Expected: []byte("\xff\xffapple\xff\xffbanana\xff\xffcherry\xff\xff"),
		},




		{
			Bytes:    []byte("apple\xff\xffbanana\xff\xffcherry"),
			Expected: []byte("apple\xff\xff\xff\xffbanana\xff\xff\xff\xffcherry"),
		},
		{
			Bytes:    []byte("\xff\xffapple\xff\xffbanana\xff\xffcherry\xff\xff"),
			Expected: []byte("\xff\xff\xff\xffapple\xff\xff\xff\xffbanana\xff\xff\xff\xffcherry\xff\xff\xff\xff"),
		},
	}

//@TODO: Add random tests.


	for testNumber, test := range tests {

		subWriter := new(bytes.Buffer)

		writer := newDataWriter(subWriter)

		n, err := writer.Write(test.Bytes)
		if nil != err {
			t.Errorf("For test #%d, did not expected an error, but actually got one: (%T) %v; for %q -> %q.", testNumber, err, err, string(test.Bytes), string(test.Expected))
			continue
		}

		if expected, actual := len(test.Bytes), n; expected != actual {
			t.Errorf("For test #%d, expected %d, but actually got %d; for %q -> %q.", testNumber, expected, actual, string(test.Bytes), string(test.Expected))
			continue
		}

		if expected, actual := string(test.Expected), subWriter.String(); expected != actual {
			t.Errorf("For test #%d, expected %q, but actually got %q; for %q -> %q.", testNumber, expected, actual, string(test.Bytes), string(test.Expected))
			continue
		}
	}
}
