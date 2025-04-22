package gedis

import (
	"time"

	redis "github.com/GedisCaching/Gedis/server"
)

// Gedis represents a Redis-like database server with operations
type Gedis struct {
	server *redis.Server
}

type Config struct {
	Address  string
	Password string
}

// NewGedis creates a new Gedis instance from a redis.Server
func NewGedis(config Config) (*Gedis, error) {
	// Create a new server with the provided config
	serverConfig := &redis.Config{
		Address:  config.Address,
		Password: config.Password,
	}

	server, err := redis.NewServer(serverConfig)
	if err != nil {
		return nil, err
	}

	return &Gedis{server: server}, nil
}

// ----------------------- SET function -----------------------

// SET function
func (g *Gedis) Set(key string, value interface{}) {
	g.server.UpdateAccessTime()
	g.server.GetDB().Set(key, value)
}

// SetWithExpiry function
func (g *Gedis) SetWithExpiry(key string, value interface{}, expiry time.Duration) {
	g.server.UpdateAccessTime()
	g.server.GetDB().SetWithExpiry(key, value, expiry)
}

// DEXPIRE function
func (g *Gedis) DEXPIRE(key string, expiry time.Duration) error {
	g.server.UpdateAccessTime()
	return g.server.GetDB().DEXPIRE(key, expiry)
}

// RENAME function
func (g *Gedis) RENAME(KeyOld, KeyNew string) error {
	g.server.UpdateAccessTime()
	return g.server.GetDB().RENAME(KeyOld, KeyNew)
}

// ----------------------- GET, DEL, KEYS Operations -----------------------

// GET function
func (g *Gedis) Get(key string) (interface{}, bool) {
	g.server.UpdateAccessTime()
	return g.server.GetDB().Get(key)
}

// GETDEL function
func (g *Gedis) GETDEL(key string) (interface{}, bool) {
	g.server.UpdateAccessTime()
	return g.server.GetDB().GETDEL(key)
}

// DEL function
func (g *Gedis) Delete(key string) bool {
	g.server.UpdateAccessTime()
	return g.server.GetDB().Delete(key)
}

// KEYS function
func (g *Gedis) Keys() []string {
	g.server.UpdateAccessTime()
	return g.server.GetDB().Keys()
}

// ----------------------- NUMERIC Operations -----------------------

// INCR function
func (g *Gedis) Incr(key string) (int, error) {
	g.server.UpdateAccessTime()
	return g.server.GetDB().Incr(key)
}

// DECR function
func (g *Gedis) Decr(key string) (int, error) {
	g.server.UpdateAccessTime()
	return g.server.GetDB().Decr(key)
}

// ----------------------- List Operations -----------------------

// LPUSH function
func (g *Gedis) LPush(key string, values ...interface{}) (int, error) {
	g.server.UpdateAccessTime()
	return g.server.GetDB().LPush(key, values...)
}

// RPUSH function
func (g *Gedis) RPush(key string, values ...interface{}) (int, error) {
	g.server.UpdateAccessTime()
	return g.server.GetDB().RPush(key, values...)
}

// LRANGE function
func (g *Gedis) LRange(key string, start, stop int) ([]interface{}, error) {
	g.server.UpdateAccessTime()
	return g.server.GetDB().LRange(key, start, stop)
}

// LPOP function
func (g *Gedis) LPop(key string) (interface{}, error) {
	g.server.UpdateAccessTime()
	return g.server.GetDB().LPop(key)
}

// RPOP function
func (g *Gedis) RPop(key string) (interface{}, error) {
	g.server.UpdateAccessTime()
	return g.server.GetDB().RPop(key)
}

// GET list length
func (g *Gedis) LLen(key string) (int, error) {
	g.server.UpdateAccessTime()
	return g.server.GetDB().LLen(key)
}

// SET list element
func (g *Gedis) LSet(key string, index int, value interface{}) error {
	g.server.UpdateAccessTime()
	return g.server.GetDB().LSet(key, index, value)
}

// ------------------------- TTL Operations -----------------------

// TTL function
func (g *Gedis) TTL(key string) (time.Duration, bool) {
	g.server.UpdateAccessTime()
	return g.server.GetDB().TTL(key)
}

// ------------------------- Sorted Set Operations -----------------------

// ZADD function
func (g *Gedis) ZAdd(key string, scoreMembers map[string]float64) int {
	g.server.UpdateAccessTime()
	return g.server.GetDB().ZADD(key, scoreMembers)
}

// ZRANGE function
func (g *Gedis) ZRange(key string, start, stop int, withScores bool) []interface{} {
	g.server.UpdateAccessTime()
	return g.server.GetDB().ZRANGE(key, start, stop, withScores)
}

// ZRANK function
func (g *Gedis) ZRank(key, member string) (int, bool) {
	g.server.UpdateAccessTime()
	return g.server.GetDB().ZRANK(key, member)
}
