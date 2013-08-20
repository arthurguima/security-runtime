// Simple GO server
package main

import (
  "net"     // provides the Listener and Conn types that hide many of the details of setting up socket connections
  "bufio"   // provides buffered read methods, simplifying common tasks like reading lines from a socket
  "strconv" // function Itoa() that converts an integer to a string
  "fmt"     // for printing strings to the console
  "bytes"   // Compare bytes
)

const PORT = 8000 // port that the server is going to listen

func main() {

  // we start by declaring and initializing a new listener for the server
  server, err := net.Listen("tcp", "10.16.1.14:" + strconv.Itoa(PORT))
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
      // sends hello!
         byteMessage := []byte("Hello! Connection established with: " + listener.Addr().String() + ".\n")
         client.Write(byteMessage) 
      // send the client, of type net.Conn to the channel 
         ch <- client
    }

   }()
   return ch
}


// controls the main connection with client
func handleConn(client net.Conn) {

  // used to manage the connection
    unsec := []byte("unsec")
    sec   := []byte("sec")
    unsec_ch := make(chan int)
    //sec_ch   make(chan )

  // buffers the client req
  b := bufio.NewReader(client)

  for {
    line, err := b.ReadBytes('\n')
    if err != nil { // EOF, or worse
      break
    }

    if(bytes.HasPrefix(line, unsec) == true){
      fmt.Print("Unsec connection with ", client.RemoteAddr(), "\n")
      go HandleUnsec(client, unsec_ch)
      //sec_ch <- 0

    } else if(bytes.HasPrefix(line, sec) == true){
      fmt.Print("Sec connection with ", client.RemoteAddr(), "\n")
      //go HandleSec()
      unsec_ch <- 1 
    }
  }
}

func HandleUnsec (client net.Conn, exit chan int ) {
  // leave to the SO the responsability to choose an avaliable port
  unsec_server, err := net.Listen("tcp", "10.16.1.14:0")
  if unsec_server == nil {
    // exits the application
    panic(err)
  }

  // announces to client the address 
  unsec_addr := []byte( unsec_server.Addr().String())
  client.Write(unsec_addr)

  //Waits for client to connect
    unsec_client, err := unsec_server.Accept()
    if unsec_client == nil {
        fmt.Printf("couldn't accept: %s\n", err)
        panic(err)
    }

  // buffers the client req
     b := bufio.NewReader(unsec_client) 

     interrupt := false

  for {
    select {
      case <- exit:
        interrupt = true
        break
      default:
        line, err := b.ReadBytes('\n')
        if err != nil { // EOF, or worse
          break
        }
        fmt.Print("Unsec Message Received: \n") 
        unsec_client.Write(line)
    }
    if (interrupt == true){
      break
    }
  }
  unsec_server.Close()
}

/*func HandleSec (client net.Conn){
}*/
