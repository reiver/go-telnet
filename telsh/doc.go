/*
Package telsh provides "middleware" (for the telnet package) that can be used to implement a TELNET server
that provides a "shell" interface.

Shell interfaces you may be familiar with include: "bash", "csh", "sh", "zsk", etc.

TELNET Server

Here is an example usage:

	package main
	
	import (
		"github.com/reiver/go-oi"
		"github.com/reiver/go-telnet"
		"github.com/reiver/go-telnet/telsh"

		"io"
	)

	func main() {
		
		telnetHandler := telsh.NewShellHandler()
		
		if err := telnetHandler.RegisterElse(
			telsh.ProducerFunc(
				func(ctx telsh.Context, name string, args ...string) telsh.Handler {
					return telsh.PromoteHandlerFunc(
						func(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser) error {
							oi.LongWrite(stdout, []byte{'w','a','t','?', '\r','\n'})
							
							return nil
						},
					)
				},
			),
		); nil != err {
			panic(err)
		}
		
		if err := telnetHandler.Register("help",
			telsh.ProducerFunc(
				func(ctx telsh.Context, name string, args ...string) telsh.Handler {
				return telsh.PromoteHandlerFunc(
						func(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser) error {
							oi.LongWrite(stdout, []byte{'r','t','f','m','!', '\r','\n'})
							
							return nil
						},
					)
				},
			),
		); nil != err {
			panic(err)
		}
		
		err := telnet.ListenAndServe(":5555", telnetHandler)
		if nil != err {
			//@TODO: Handle this error better.
			panic(err)
		}
	}

Here is a more "unpacked" example:

	package main
	
	
	import (
		"github.com/reiver/go-oi"
		"github.com/reiver/go-telnet"
		"github.com/reiver/go-telnet/telsh"
		
		"fmt"
		"io"
		"time"
	)
	
	
	var (
		shellHandler := telsh.NewShellHandler()
	)
	
	
	func init() {
		
		shellHandler.Register("dance", telsh.ProducerFunc(producer))
		
		
		shellHandler.WelcomeMessage = `
	 __          __ ______  _        _____   ____   __  __  ______ 
	 \ \        / /|  ____|| |      / ____| / __ \ |  \/  ||  ____|
	  \ \  /\  / / | |__   | |     | |     | |  | || \  / || |__   
	   \ \/  \/ /  |  __|  | |     | |     | |  | || |\/| ||  __|  
	    \  /\  /   | |____ | |____ | |____ | |__| || |  | || |____ 
	     \/  \/    |______||______| \_____| \____/ |_|  |_||______|
	
	`
	}
	
	
	func producer(ctx telsh.Context, name string, args ...string) telsh.Handler{
		return telsh.PromoteHandlerFunc(handler)
	}
	
	
	func handler(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser) error {
		for i:=0; i<20; i++ {
			oi.LongWriteString(stdout, "\r⠋")
			time.Sleep(50*time.Millisecond)
			
			oi.LongWriteString(stdout, "\r⠙")
			time.Sleep(50*time.Millisecond)
			
			oi.LongWriteString(stdout, "\r⠹")
			time.Sleep(50*time.Millisecond)
			
			oi.LongWriteString(stdout, "\r⠸")
			time.Sleep(50*time.Millisecond)
			
			oi.LongWriteString(stdout, "\r⠼")
			time.Sleep(50*time.Millisecond)
			
			oi.LongWriteString(stdout, "\r⠴")
			time.Sleep(50*time.Millisecond)
			
			oi.LongWriteString(stdout, "\r⠦")
			time.Sleep(50*time.Millisecond)
			
			oi.LongWriteString(stdout, "\r⠧")
			time.Sleep(50*time.Millisecond)
			
			oi.LongWriteString(stdout, "\r⠇")
			time.Sleep(50*time.Millisecond)
			
			oi.LongWriteString(stdout, "\r⠏")
			time.Sleep(50*time.Millisecond)
		}
		oi.LongWriteString(stdout, "\r \r\n")

		return nil
	}
	
	
	func main() {
		
		addr := ":5555"
		if err := telnet.ListenAndServe(addr, shellHandler); nil != err {
			panic(err)
		}
	}

*/
package telsh
