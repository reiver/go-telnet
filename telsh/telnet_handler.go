package telsh


import (
	"github.com/reiver/go-oi"
	"github.com/reiver/go-telnet"

	"bytes"
	"io"
	"sync"
)


const (
	defaultExitCommandName = "exit"
	defaultPrompt          = "$ "
	defaultWelcomeMessage  = "Welcome!\r\n"
	defaultExitMessage     = "Goodbye!\r\n"
)


type ShellHandler struct {
	muxtex sync.RWMutex
	shellCommands map[string]Producer
	elseHandler Handler

	ExitCommandName string
	Prompt          string
	WelcomeMessage  string
	ExitMessage     string
}


func NewShellHandler() *ShellHandler {
	shellCommands := map[string]Producer{}

	handler := ShellHandler{
		shellCommands:shellCommands,

		Prompt:          defaultPrompt,
		ExitCommandName: defaultExitCommandName,
		WelcomeMessage:  defaultWelcomeMessage,
		ExitMessage:     defaultExitMessage,
	}

	return &handler
}


func (handler *ShellHandler) Register(name string, producer Producer) error {

	handler.muxtex.Lock()
	handler.shellCommands[name] = producer
	handler.muxtex.Unlock()

	return nil
}


func (handler *ShellHandler) Else(elseHandler Handler) error {

	handler.muxtex.Lock()
	handler.elseHandler = elseHandler
	handler.muxtex.Unlock()

	return nil
}


func (handler *ShellHandler) ServeTELNET(ctx telnet.Context, w io.Writer, r io.Reader) {

	colonSpaceCommandNotFoundEL := []byte(": command not found\r\n")


	writer := telnet.NewDataWriter(w)
	reader := telnet.NewDataReader(r)


	var prompt          bytes.Buffer
	var exitCommandName bytes.Buffer
	var welcomeMessage  bytes.Buffer
	var exitMessage     bytes.Buffer

	prompt.WriteString(handler.Prompt)
	exitCommandName.WriteString(handler.ExitCommandName)
	welcomeMessage.WriteString(handler.WelcomeMessage)
	exitMessage.WriteString(handler.ExitMessage)

	promptBytes          := prompt.Bytes()
	exitCommandNameBytes := exitCommandName.Bytes()
	welcomeMessageBytes  := welcomeMessage.Bytes()
	exitMessageBytes     := exitMessage.Bytes()


//@TODO: Should we check for (potential) errors coming from these LongWrites?
	oi.LongWrite(writer, welcomeMessageBytes)
	oi.LongWrite(writer, promptBytes)


	var buffer [1]byte // Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
	p := buffer[:]

	var line bytes.Buffer

	for {
		// Read 1 byte.
		n, err := reader.Read(p)
		if n <= 0 && nil == err {
			continue
		} else if n <= 0 && nil != err {
			break
		}


		line.WriteByte(p[0])


		if '\n' == p[0] {
			lineBytes := line.Bytes()

			if 2 == len(lineBytes) && '\r' == lineBytes[0] && '\n' == lineBytes[1] {
				line.Reset()
				oi.LongWrite(writer, promptBytes)
				continue
			}


			fields := bytes.Fields(lineBytes)
			field0 := fields[0]

			if 0 == bytes.Compare(exitCommandNameBytes, field0) {
				oi.LongWrite(writer, exitMessageBytes)
				return
			}


			oi.LongWrite(writer, field0)
			oi.LongWrite(writer, colonSpaceCommandNotFoundEL)
			oi.LongWrite(writer, promptBytes)
			line.Reset()
		}


//@TODO: Are there any special errors we should be dealing with separately?
		if nil != err {
			break
		}
	}


	oi.LongWrite(writer, exitMessageBytes)
	return
}
