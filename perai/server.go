// Simple GO server
package main

import (
  "net"     // provides the Listener and Conn types that hide many of the details of setting up socket connections
  "bufio"   // provides buffered read methods, simplifying common tasks like reading lines from a socket
  "strconv" // function Itoa() that converts an integer to a string
  "fmt"     // for printing strings to the console
)

const PORT = 8000 // port that the server is going to listen

func main() {

  // we start by declaring and initializing a new listener for the server
  server, err := net.Listen("tcp", ":" + strconv.Itoa(PORT))
  if server == nil {
    // exits the application
    panic(err)
  }

  conns := clientConns(server) // this is the channel we’ll use for getting new client connections.
  // infinite loop
  // each time we start a goroutine, with the value receive operation on our client connections channel.
  // the unary operator <- blocks until a value is available on the channel (a new client having connected)
  for { 
    go handleConn(<-conns)
  }
}

func clientConns(listener net.Listener) chan net.Conn {

   // channel that corresponds to the type that we’ll be got from calling Accept() on listener connection object
   ch := make(chan net.Conn)

   // anonymous goroutine which runs in an infinite loop, constantly accepting new connections  
   go func(){

    for{
      // blocks as long as there are no new clients to deal with 
      client, err := listener.Accept()
      if client == nil {
        fmt.Printf("couldn't accept: %s\n", err)
        continue
      }
      fmt.Printf("New connection with: %v established\n", client.RemoteAddr())
      // send the client, of type net.Conn to the channel 
      ch <- client
    }

   }()
   return ch
}

func handleConn(client net.Conn) {

  //Sends Hello!
  byteMessage := []byte("Hello! Connection established with " + "\n")
  client.Write(byteMessage)

  b := bufio.NewReader(client)
  for {
    line, err := b.ReadBytes('\n')
    if err != nil { // EOF, or worse
      break
    }
  client.Write(line)
  }
}

