package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to Redis clone. Type commands (e.g., SET key value, GET key, DEL key). Type 'EXIT' to quit.")

	reader := bufio.NewReader(os.Stdin)

	for {
		// Read input from the user
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if strings.ToUpper(input) == "EXIT" {
			fmt.Println("Exiting...")
			break
		}

		// Send the input to the server
		_, err = conn.Write([]byte(input + "\n"))
		if err != nil {
			fmt.Println("Error writing to server:", err)
			break
		}

		// Read and print the response from the server
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from server:", err)
			break
		}

		fmt.Println(strings.TrimSpace(response))
	}
}
