# go-telnet

Package **telnet** provides TELNET client and server implementations, for Go programming language,
in a style similar to the "net/http" library that is part of the Go standard library,
including support for "middleware".

## Documention

Online documentation, which includes examples, can be found at: http://godoc.org/github.com/reiver/go-telnet

[![GoDoc](https://godoc.org/github.com/reiver/go-telnet?status.svg)](https://godoc.org/github.com/reiver/go-telnet)

## TELNET Server Example
```
package main

import (
	"github.com/reiver/go-telnet"
)

func main() {

	var handler telnet.Handler = telnet.EchoHandler
	
	err := telnet.ListenAndServe(":5555", handler)
	if nil != err {
		//@TODO: Handle this error better.
		panic(err)
	}
}

```
