package telsh

import (
	"github.com/matjam/go-telnet"

	"bytes"
	"strings"

	"testing"
)

func TestServeTELNETCommandNotFound(t *testing.T) {

	tests := []struct {
		ClientSends string
		Expected    string
	}{
		{
			ClientSends: "\r\n",
			Expected:    "",
		},

		{
			ClientSends: "apple\r\n",
			Expected:    "apple: command not found\r\n",
		},
		{
			ClientSends: "banana\r\n",
			Expected:    "banana: command not found\r\n",
		},
		{
			ClientSends: "cherry\r\n",
			Expected:    "cherry: command not found\r\n",
		},

		{
			ClientSends: "\t\r\n",
			Expected:    "",
		},
		{
			ClientSends: "\t\t\r\n",
			Expected:    "",
		},
		{
			ClientSends: "\t\t\t\r\n",
			Expected:    "",
		},

		{
			ClientSends: " \r\n",
			Expected:    "",
		},
		{
			ClientSends: "  \r\n",
			Expected:    "",
		},
		{
			ClientSends: "   \r\n",
			Expected:    "",
		},

		{
			ClientSends: " \t\r\n",
			Expected:    "",
		},
		{
			ClientSends: "\t \r\n",
			Expected:    "",
		},

		{
			ClientSends: "ls -alF\r\n",
			Expected:    "ls: command not found\r\n",
		},
	}

	for testNumber, test := range tests {

		shellHandler := NewShellHandler()
		if nil == shellHandler {
			t.Errorf("For test #%d, did not expect to get nil, but actually got it: %v; for client sent: %q", testNumber, shellHandler, test.ClientSends)
			continue
		}

		ctx := telnet.NewContext()

		var buffer bytes.Buffer

		shellHandler.ServeTELNET(ctx, &buffer, strings.NewReader(test.ClientSends))

		if expected, actual := shellHandler.WelcomeMessage+shellHandler.Prompt+test.Expected+shellHandler.Prompt+shellHandler.ExitMessage, buffer.String(); expected != actual {
			t.Errorf("For test #%d, expect %q, but actually got %q; for client sent: %q", testNumber, expected, actual, test.ClientSends)
			continue
		}
	}
}
