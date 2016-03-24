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


## TELNET Story

The TELNET protocol is best known for providing a means of connecting to a remote computer,
using a (text-based) *shell* interface, and being able to interact with it, (more or less)
as if you were sitting at that computer.

(*Shells* are also known as *command-line interfaces* or *CLIs*.)

Although this was the original usage of the TELNET protocol, **it can be (and is) used for
other purposes as well**.

The TELNET protocol came from an era in computing when text-based *shell* interface where the
common way of interacting with computers.

The common interface for computers during this era was a keyboard and a monochromatic (i.e., single color) text-based
monitors (called "video terminals").

(The word "video" in that era of computing did not refer to things such as *movies*. But instead
was meant to contrast it with paper. In particular, the *teletype* machines, which were typewriter
like devices that had a keyboard, but instead of having a monitor have paper that way printed onto.)

