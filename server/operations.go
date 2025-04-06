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
func (s *Server) Keys() []interface{} {
	return s.db.Keys()
}

// INCR function
func (s *Server) Incr(key interface{}) (int, error) {
	return s.db.Incr(key)
}

// DECR function
func (s *Server) Decr(key interface{}) (int, error) {
	return s.db.Decr(key)
}
