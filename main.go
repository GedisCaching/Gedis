package main

import (
	"fmt"
	"github.com/GedisCaching/Gedis/RESP"
	"io"
	"net"
	"os"
)

const defaultAddress = "0.0.0.0:7000"

func main() {
	l, err := net.Listen("tcp", defaultAddress)
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Starting server at", defaultAddress)

	// Listen for inputs and respond
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2014) // store out stuff somewhere

	for {
		len, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading: %#v\n", err)
			}
			break
		}

		// Parse the command out
		command := buf[:len]
		response := RESP.Parse(command)

		// Write the response back to the connection
		_, Reserr := conn.Write([]byte(fmt.Sprintf("%v\r\n", response)))
		if Reserr != nil {
			fmt.Printf("Error reading: %#v\n", Reserr)
			break
		}
	}
}
