package telsh


import (
	"github.com/reiver/go-oi"
	"github.com/reiver/go-telnet"

	"bytes"
	"io"
	"strings"
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
	producers map[string]Producer
	elseProducer Producer

	ExitCommandName string
	Prompt          string
	WelcomeMessage  string
	ExitMessage     string
}


func NewShellHandler() *ShellHandler {
	producers := map[string]Producer{}

	telnetHandler := ShellHandler{
		producers:producers,

		Prompt:          defaultPrompt,
		ExitCommandName: defaultExitCommandName,
		WelcomeMessage:  defaultWelcomeMessage,
		ExitMessage:     defaultExitMessage,
	}

	return &telnetHandler
}


func (telnetHandler *ShellHandler) Register(name string, producer Producer) error {

	telnetHandler.muxtex.Lock()
	telnetHandler.producers[name] = producer
	telnetHandler.muxtex.Unlock()

	return nil
}

func (telnetHandler *ShellHandler) MustRegister(name string, producer Producer) *ShellHandler {
	if err := telnetHandler.Register(name, producer); nil != err {
		panic(err)
	}

	return telnetHandler
}


func (telnetHandler *ShellHandler) RegisterElse(producer Producer) error {

	telnetHandler.muxtex.Lock()
	telnetHandler.elseProducer = producer
	telnetHandler.muxtex.Unlock()

	return nil
}

func (telnetHandler *ShellHandler) MustRegisterElse(producer Producer) *ShellHandler {
	if err := telnetHandler.RegisterElse(producer); nil != err {
			panic(err)
	}

	return telnetHandler
}


func (telnetHandler *ShellHandler) ServeTELNET(ctx telnet.Context, w io.Writer, r io.Reader) {

	colonSpaceCommandNotFoundEL := []byte(": command not found\r\n")


	writer := telnet.NewDataWriter(w)
	reader := telnet.NewDataReader(r)


	var prompt          bytes.Buffer
	var exitCommandName string
	var welcomeMessage  bytes.Buffer
	var exitMessage     bytes.Buffer

	prompt.WriteString(telnetHandler.Prompt)
	welcomeMessage.WriteString(telnetHandler.WelcomeMessage)
	exitMessage.WriteString(telnetHandler.ExitMessage)

	promptBytes          := prompt.Bytes()
	welcomeMessageBytes  := welcomeMessage.Bytes()
	exitMessageBytes     := exitMessage.Bytes()

	exitCommandName = telnetHandler.ExitCommandName


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
			lineString := line.String()

			if "\r\n" == lineString {
				line.Reset()
				oi.LongWrite(writer, promptBytes)
				continue
			}


//@TODO: support piping.
			fields := strings.Fields(lineString)
			field0 := fields[0]

			if exitCommandName == field0 {
				oi.LongWrite(writer, exitMessageBytes)
				return
			}


			var producer Producer

			telnetHandler.muxtex.RLock()
			var ok bool
			producer, ok = telnetHandler.producers[field0]
			telnetHandler.muxtex.RUnlock()

			if !ok {
				telnetHandler.muxtex.RLock()
				producer = telnetHandler.elseProducer
				telnetHandler.muxtex.RUnlock()
			}

			if nil == producer {
//@TODO: Don't convert that to []byte! think this creates "garbage" (for collector).
				oi.LongWrite(writer, []byte(field0))
				oi.LongWrite(writer, colonSpaceCommandNotFoundEL)
				line.Reset()
				oi.LongWrite(writer, promptBytes)
				continue
			}

			handler := producer.Produce(ctx, field0, fields[1:]...)
			if nil == handler {
				oi.LongWrite(writer, []byte(field0))
//@TODO: Need to use a different error message.
				oi.LongWrite(writer, colonSpaceCommandNotFoundEL)
				line.Reset()
				oi.LongWrite(writer, promptBytes)
				continue
			}

//@TODO: Wire up the stdin, stdout, stderr of the handler.

			if stdoutPipe, err := handler.StdoutPipe(); nil != err {
//@TODO:                              
			} else if nil == stdoutPipe {
//@TODO:                              
			} else {
				go connect(writer, stdoutPipe)
			}


			if err := handler.Run(); nil != err {
//@TODO:                                    
			}
			line.Reset()
			oi.LongWrite(writer, promptBytes)
		}


//@TODO: Are there any special errors we should be dealing with separately?
		if nil != err {
			break
		}
	}


	oi.LongWrite(writer, exitMessageBytes)
	return
}



func connect(writer io.Writer, reader io.Reader) {

	var buffer [1]byte // Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
	p := buffer[:]

	for {
		// Read 1 byte.
		n, err := reader.Read(p)
		if n <= 0 && nil == err {
			continue
		} else if n <= 0 && nil != err {
			break
		}

//@TODO: Should we be checking for errors?
		oi.LongWrite(writer, p)
	}
}
