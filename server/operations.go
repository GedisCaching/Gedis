package redis

import (
	"time"
)

// ----------------------- DB Operations -----------------------

// SET function
func (s *Server) Set(key string, value interface{}) {
	s.db.Set(key, value)
}

func (s *Server) SetWithExpiry(key string, value interface{}, expiry int, EX, PX bool) {
	// EX: Set expiry in seconds
	if EX {
		s.db.SetWithExpiry(key, value, time.Duration(expiry)*time.Second)
	}
	// PX: Set expiry in milliseconds
	if PX {
		s.db.SetWithExpiry(key, value, time.Duration(expiry)*time.Millisecond)
	}
}

// GET function
func (s *Server) Get(key string) (interface{}, bool) {
	return s.db.Get(key)
}

// DEL function
func (s *Server) Delete(key string) bool {
	return s.db.Delete(key)
}

// KEYS function
func (s *Server) Keys() []string {
	return s.db.Keys()
}

// INCR function
func (s *Server) Incr(key string) (int, error) {
	return s.db.Incr(key)
}

// DECR function
func (s *Server) Decr(key string) (int, error) {
	return s.db.Decr(key)
}

// ----------------------- List Operations -----------------------

// LPUSH function
func (s *Server) LPush(key string, values ...interface{}) (int, error) {
	return s.db.LPush(key, values...)
}

// RPUSH function
func (s *Server) RPush(key string, values ...interface{}) (int, error) {
	return s.db.RPush(key, values...)
}

// LRANGE function
func (s *Server) LRange(key string, start, stop int) ([]interface{}, error) {
	return s.db.LRange(key, start, stop)
}
