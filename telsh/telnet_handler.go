package telsh

import (
	"bufio"

	"github.com/reiver/go-oi"
	"github.com/reiver/go-telnet"

	"bytes"
	"io"
	"strings"
	"sync"
)

const (
	defaultExitCommandName = "exit"
	defaultPrompt          = "§ "
	defaultWelcomeMessage  = "\r\nWelcome!\r\n"
	defaultExitMessage     = "\r\nGoodbye!\r\n"
)

type ShellHandler struct {
	muxtex       sync.RWMutex
	producers    map[string]Producer
	elseProducer Producer

	ExitCommandName string
	Prompt          string
	WelcomeMessage  string
	ExitMessage     string
}

func NewShellHandler() *ShellHandler {
	producers := map[string]Producer{}

	telnetHandler := ShellHandler{
		producers: producers,

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

func (telnetHandler *ShellHandler) RegisterHandlerFunc(name string, handlerFunc HandlerFunc) error {

	produce := func(ctx telnet.Context, name string, args ...string) Handler {
		return PromoteHandlerFunc(handlerFunc, args...)
	}

	producer := ProducerFunc(produce)

	return telnetHandler.Register(name, producer)
}

func (telnetHandler *ShellHandler) MustRegisterHandlerFunc(name string, handlerFunc HandlerFunc) *ShellHandler {
	if err := telnetHandler.RegisterHandlerFunc(name, handlerFunc); nil != err {
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

func (telnetHandler *ShellHandler) ServeTELNET(ctx telnet.Context, writer telnet.Writer, reader telnet.Reader) {

	logger := ctx.Logger()
	if nil == logger {
		logger = internalDiscardLogger{}
	}

	colonSpaceCommandNotFoundEL := []byte(": command not found\r\n")

	var prompt bytes.Buffer
	var exitCommandName string
	var welcomeMessage string
	var exitMessage string

	prompt.WriteString(telnetHandler.Prompt)

	promptBytes := prompt.Bytes()

	exitCommandName = telnetHandler.ExitCommandName
	welcomeMessage = telnetHandler.WelcomeMessage
	exitMessage = telnetHandler.ExitMessage

	if _, err := oi.LongWriteString(writer, welcomeMessage); nil != err {
		logger.Errorf("Problem long writing welcome message: %v", err)
		return
	}
	logger.Debugf("Wrote welcome message: %q.", welcomeMessage)
	if _, err := oi.LongWrite(writer, promptBytes); nil != err {
		logger.Errorf("Problem long writing prompt: %v", err)
		return
	}
	logger.Debugf("Wrote prompt: %q.", promptBytes)

	buffered := bufio.NewReader(reader)

	var err error
	var line string

	for err == nil {
		line, err = buffered.ReadString('\n')
		if err != nil {
			break
		}

		if "\r\n" == line {
			_, err = oi.LongWrite(writer, promptBytes)
			continue
		}

		//@TODO: support piping.
		fields := strings.Fields(line)
		logger.Debugf("Have %d tokens.", len(fields))
		logger.Tracef("Tokens: %v", fields)
		if len(fields) <= 0 {
			_, err = oi.LongWrite(writer, promptBytes)
			continue
		}

		field0 := fields[0]

		if exitCommandName == field0 {
			break
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
			_, err = oi.LongWrite(writer, promptBytes)
			continue
		}

		handler := producer.Produce(ctx, field0, fields[1:]...)
		if nil == handler {
			oi.LongWrite(writer, []byte(field0))
			//@TODO: Need to use a different error message.
			oi.LongWrite(writer, colonSpaceCommandNotFoundEL)
			_, err = oi.LongWrite(writer, promptBytes)
			continue
		}

		//@TODO: Wire up the stdin, stdout, stderr of the handler.

		if stdoutPipe, err := handler.StdoutPipe(); nil != err {
			//@TODO:
		} else if nil == stdoutPipe {
			//@TODO:
		} else {
			connect(ctx, writer, stdoutPipe)
		}

		if stderrPipe, err := handler.StderrPipe(); nil != err {
			//@TODO:
		} else if nil == stderrPipe {
			//@TODO:
		} else {
			connect(ctx, writer, stderrPipe)
		}

		if err = handler.Run(); nil != err {
			//@TODO:
		}
		_, err = oi.LongWrite(writer, promptBytes)
	}

	oi.LongWriteString(writer, exitMessage)
	return
}

func connect(ctx telnet.Context, writer io.Writer, reader io.Reader) {

	logger := ctx.Logger()

	go func(logger telnet.Logger) {
		io.Copy(writer, reader)
	}(logger)
}
