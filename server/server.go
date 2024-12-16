package server

import (
	"bufio"
	"fmt"
	"net"
	"redis-go/internal/kvstore"
	"strings"
)

func HandleConnection(conn net.Conn, s *kvstore.Store) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		input := scanner.Text()
		parts := strings.Split(input, " ")

		if len(parts) < 2 {
			fmt.Fprintln(conn, "ERROR: Unknown command")
			continue
		}

		command := parts[0]
		args := parts[1:]

		response := s.HandleCommand(command, args)

		fmt.Fprintln(conn, response)
	}
}
