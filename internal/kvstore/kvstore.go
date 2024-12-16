package kvstore

import (
	"fmt"
	"strings"
	"sync"
)

// Store represents a thread-safe key-value store.
type Store struct {
	data map[string]string
	mu   sync.RWMutex
}

// NewStore initializes and returns a new Store.
func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

// Set stores a key-value pair in the store.
func (s *Store) Set(key, value string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
	return fmt.Sprintf("OK: %s=%s", key, value)
}

// Get retrieves the value for a given key. Returns an error message if the key does not exist.
func (s *Store) Get(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if value, exists := s.data[key]; exists {
		return value
	}
	return "ERROR: Key not found"
}

// Delete removes a key-value pair from the store. Returns a confirmation message.
func (s *Store) Delete(key string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; exists {
		delete(s.data, key)
		return "OK: Key deleted"
	}
	return "ERROR: Key not found"
}

// HandleCommand processes commands like SET, GET, and DEL with thread safety.
func (s *Store) HandleCommand(command string, args []string) string {
	if len(args) == 0 && (command == "SET" || command == "GET" || command == "DEL") {
		return "ERROR: Missing arguments"
	}

	switch strings.ToUpper(command) {
	case "SET":
		if len(args) < 2 {
			return "ERROR: SET requires a key and value"
		}
		return s.Set(args[0], strings.Join(args[1:], " "))
	case "GET":
		return s.Get(args[0])
	case "DEL":
		return s.Delete(args[0])
	default:
		return "ERROR: Unknown command"
	}
}
