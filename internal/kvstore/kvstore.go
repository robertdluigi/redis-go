package kvstore

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type Store struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewStore() *Store {
	return &Store{data: make(map[string]string)}
}

func (s *Store) Set(key, value string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	return "OK"
}

func (s *Store) Get(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if value, exists := s.data[key]; exists {
		return value
	}
	return "ERROR: Key not found"
}

func (s *Store) Delete(key string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.data[key]; exists {
		delete(s.data, key)
		return "OK: Key deleted"
	}
	return "ERROR: Key not found"
}

func (s *Store) AdjustBy(key string, amount int) string {
	value, err := strconv.Atoi(s.Get(key))
	if err != nil {
		return "ERROR: Value is not a number"
	}
	newValue := value + amount
	s.Set(key, strconv.Itoa(newValue))
	return fmt.Sprintf("OK: %d", newValue)
}

func (s *Store) HandleCommand(command string, args []string) string {
	if len(args) == 0 && (command == "SET" || command == "GET" || command == "DEL") {
		return "ERROR: Missing arguments"
	}

	commandHandlers := map[string]func([]string) string{
		"SET":    func(args []string) string { return s.Set(args[0], strings.Join(args[1:], " ")) },
		"GET":    func(args []string) string { return s.Get(args[0]) },
		"DEL":    func(args []string) string { return s.Delete(args[0]) },
		"INCR":   func(args []string) string { return s.AdjustBy(args[0], 1) },
		"INCRBY": func(args []string) string { return s.AdjustBy(args[0], atoi(args[1])) },
		"DECR":   func(args []string) string { return s.AdjustBy(args[0], -1) },
		"DECRBY": func(args []string) string { return s.AdjustBy(args[0], -atoi(args[1])) },
	}

	handler, exists := commandHandlers[strings.ToUpper(command)]
	if !exists {
		return "ERROR: Unknown command"
	}
	return handler(args)
}

func atoi(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}
