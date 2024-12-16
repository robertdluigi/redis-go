package kvstore

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type Store struct {
	mu   sync.Mutex
	data map[string]string
	list map[string][]string
	sets map[string]map[string]bool
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
		list: make(map[string][]string),
		sets: make(map[string]map[string]bool),
	}
}

// Key-Value commands
func (s *Store) Set(key, value string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
	return "OK"
}

func (s *Store) Get(key string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if value, exists := s.data[key]; exists {
		return value
	}
	return "(nil)"
}

func (s *Store) Delete(key string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; exists {
		delete(s.data, key)
		return "(integer) 1"
	}
	return "(integer) 0"
}

// List commands
func (s *Store) LPush(key string, values ...string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.list[key]; !exists {
		s.list[key] = []string{}
	}

	s.list[key] = append(values, s.list[key]...)
	return fmt.Sprintf("(integer) %d", len(s.list[key]))
}

func (s *Store) RPush(key string, values ...string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.list[key]; !exists {
		s.list[key] = []string{}
	}

	s.list[key] = append(s.list[key], values...)
	return fmt.Sprintf("(integer) %d", len(s.list[key]))
}

func (s *Store) LRange(key string, start, end int) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if list, exists := s.list[key]; exists {
		if start < 0 {
			start = len(list) + start
		}
		if end < 0 {
			end = len(list) + end
		}
		if start < 0 {
			start = 0
		}
		if end >= len(list) {
			end = len(list) - 1
		}
		if start > end || start >= len(list) {
			return "(empty list)"
		}
		return strings.Join(list[start:end+1], ", ")
	}
	return "(nil)"
}

func (s *Store) LPop(key string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if list, exists := s.list[key]; exists && len(list) > 0 {
		value := list[0]
		s.list[key] = list[1:]
		return value
	}
	return "(nil)"
}

func (s *Store) RPop(key string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if list, exists := s.list[key]; exists && len(list) > 0 {
		value := list[len(list)-1]
		s.list[key] = list[:len(list)-1]
		return value
	}
	return "(nil)"
}

// Set commands
func (s *Store) SAdd(key string, members ...string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.sets[key]; !exists {
		s.sets[key] = make(map[string]bool)
	}

	added := 0
	for _, member := range members {
		if !s.sets[key][member] {
			s.sets[key][member] = true
			added++
		}
	}
	return fmt.Sprintf("(integer) %d", added)
}

func (s *Store) SMembers(key string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if set, exists := s.sets[key]; exists {
		members := []string{}
		for member := range set {
			members = append(members, member)
		}
		return strings.Join(members, ", ")
	}
	return "(nil)"
}

func (s *Store) SIsMember(key, member string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if set, exists := s.sets[key]; exists {
		if set[member] {
			return "(integer) 1"
		}
	}
	return "(integer) 0"
}

func (s *Store) SRem(key string, members ...string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if set, exists := s.sets[key]; exists {
		removed := 0
		for _, member := range members {
			if set[member] {
				delete(set, member)
				removed++
			}
		}
		return fmt.Sprintf("(integer) %d", removed)
	}
	return "(integer) 0"
}

// INCR/DECR commands
func (s *Store) AdjustBy(key string, amount int) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentValue, exists := s.data[key]
	if !exists {
		s.data[key] = strconv.Itoa(amount)
		return fmt.Sprintf("OK: %d", amount)
	}

	value, err := strconv.Atoi(currentValue)
	if err != nil {
		return "ERROR: Value is not a number"
	}

	newValue := value + amount
	s.data[key] = strconv.Itoa(newValue)
	return fmt.Sprintf("OK: %d", newValue)
}

func (s *Store) INCR(key string) string {
	return s.AdjustBy(key, 1)
}

func (s *Store) INCRBY(key string, amount int) string {
	return s.AdjustBy(key, amount)
}

func (s *Store) DECR(key string) string {
	return s.AdjustBy(key, -1)
}

func (s *Store) DECRBY(key string, amount int) string {
	return s.AdjustBy(key, -amount)
}

// Command handler
func (s *Store) HandleCommand(command string, args []string) string {
	commandHandlers := map[string]func([]string) string{
		"SET":    func(args []string) string { return s.Set(args[0], strings.Join(args[1:], " ")) },
		"GET":    func(args []string) string { return s.Get(args[0]) },
		"DEL":    func(args []string) string { return s.Delete(args[0]) },
		"LPUSH":  func(args []string) string { return s.LPush(args[0], args[1:]...) },
		"RPUSH":  func(args []string) string { return s.RPush(args[0], args[1:]...) },
		"LRANGE": func(args []string) string { return s.LRange(args[0], atoi(args[1]), atoi(args[2])) },
		"LPOP":   func(args []string) string { return s.LPop(args[0]) },
		"RPOP":   func(args []string) string { return s.RPop(args[0]) },
		"SADD":   func(args []string) string { return s.SAdd(args[0], args[1:]...) },
		"SMEMBERS": func(args []string) string {
			return s.SMembers(args[0])
		},
		"SISMEMBER": func(args []string) string { return s.SIsMember(args[0], args[1]) },
		"SREM":      func(args []string) string { return s.SRem(args[0], args[1:]...) },
		"INCR":      func(args []string) string { return s.INCR(args[0]) },
		"INCRBY":    func(args []string) string { return s.INCRBY(args[0], atoi(args[1])) },
		"DECR":      func(args []string) string { return s.DECR(args[0]) },
		"DECRBY":    func(args []string) string { return s.DECRBY(args[0], atoi(args[1])) },
	}

	if handler, exists := commandHandlers[command]; exists {
		return handler(args)
	}
	return "ERROR: Unknown command"
}

func atoi(s string) int {
	value, _ := strconv.Atoi(s)
	return value
}
