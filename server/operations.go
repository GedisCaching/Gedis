package redis

import (
	"time"
)

// ----------------------- SET function -----------------------

// SET function
func (s *Server) Set(key string, value interface{}) {
	s.db.Set(key, value)
}

// SetWithExpiry function
func (s *Server) SetWithExpiry(key string, value interface{}, expiry time.Duration) {
	s.db.SetWithExpiry(key, value, expiry)
}

// DEXPIRE function
func (s *Server) DEXPIRE(key string, expiry time.Duration) error {
	return s.db.DEXPIRE(key, expiry)
}

// RENAME function
func (s *Server) RENAME(KeyOld, KeyNew string) error {
	return s.db.RENAME(KeyOld, KeyNew)
}

// ----------------------- GET, DEL, KEYS Operations -----------------------

// GET function
func (s *Server) Get(key string) (interface{}, bool) {
	return s.db.Get(key)
}

// GETDEL function
func (s *Server) GETDEL(key string) (interface{}, bool) {
	return s.db.GETDEL(key)
}

// DEL function
func (s *Server) Delete(key string) bool {
	return s.db.Delete(key)
}

// KEYS function
func (s *Server) Keys() []string {
	return s.db.Keys()
}

// ----------------------- NUMERIC Operations -----------------------

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

// LPOP function
func (s *Server) LPop(key string) (interface{}, error) {
	return s.db.LPop(key)
}

// RPOP function
func (s *Server) RPop(key string) (interface{}, error) {
	return s.db.RPop(key)
}

// GET list length
func (s *Server) LLen(key string) (int, error) {
	return s.db.LLen(key)
}

// SET list element
func (s *Server) LSet(key string, index int, value interface{}) error {
	return s.db.LSet(key, index, value)
}

// ------------------------- TTL Operations -----------------------

// TTL function
func (s *Server) TTL(key string) (time.Duration, bool) {
	return s.db.TTL(key)
}
