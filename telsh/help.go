package telsh


import (
	"github.com/reiver/go-oi"

	"io"
	"sort"
)


type internalHelpProducer struct {
	shellHandler *ShellHandler
}


func Help(shellHandler *ShellHandler) Producer {
	producer := internalHelpProducer{
		shellHandler:shellHandler,
	}

	return &producer
}


func (producer *internalHelpProducer) Produce(Context, string, ...string) Handler {
	return newHelpHandler(producer)
}


type internalHelpHandler struct {
	helpProducer *internalHelpProducer

	err error

	stdin  io.ReadCloser
	stdout io.WriteCloser
	stderr io.WriteCloser

	stdinPipe  io.WriteCloser
	stdoutPipe io.ReadCloser
	stderrPipe io.ReadCloser
}


func newHelpHandler(helpProducer *internalHelpProducer) *internalHelpHandler {
	stdin,      stdinPipe := io.Pipe()
	stdoutPipe, stdout    := io.Pipe()
	stderrPipe, stderr    := io.Pipe()

	handler := internalHelpHandler{
		helpProducer:helpProducer,

		err:nil,

		stdin:stdin,
		stdout:stdout,
		stderr:stderr,

		stdinPipe:stdinPipe,
		stdoutPipe:stdoutPipe,
		stderrPipe:stderrPipe,
	}

	return &handler
}




func (handler *internalHelpHandler) Run() error {
	if nil != handler.err {
		return handler.err
	}

	//@TODO: Should this be reaching inside of ShellHandler? Maybe there should be ShellHandler public methods instead.
	keys := make([]string, 1+len(handler.helpProducer.shellHandler.producers))
	i:=0
	for key,_ := range handler.helpProducer.shellHandler.producers {
		keys[i] = key
		i++
	}
	keys[i] = handler.helpProducer.shellHandler.ExitCommandName
	sort.Strings(keys)
	for _, key := range keys {
		oi.LongWriteString(handler.stdout, key)
		oi.LongWriteString(handler.stdout, "\r\n")
	}

	handler.stdin.Close()
	handler.stdout.Close()
	handler.stderr.Close()

	return handler.err
}

func (handler *internalHelpHandler) StdinPipe() (io.WriteCloser, error) {
	if nil != handler.err {
		return nil, handler.err
	}

	return handler.stdinPipe, nil
}

func (handler *internalHelpHandler) StdoutPipe() (io.ReadCloser, error) {
	if nil != handler.err {
		return nil, handler.err
	}

	return handler.stdoutPipe, nil
}

func (handler *internalHelpHandler) StderrPipe() (io.ReadCloser, error) {
	if nil != handler.err {
		return nil, handler.err
	}

	return handler.stderrPipe, nil
}
