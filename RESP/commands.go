package RESP

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	responses "github.com/GedisCaching/Gedis/responses"
)

type StoreItem struct {
	Value  string
	Expiry time.Time
	Mutex  sync.Mutex
}

// Example usage of the store to avoid unused variable error
var store = make(map[string]*StoreItem)

func SET(key string, value string) {
	store[key] = &StoreItem{
		Value: value,
	}
}

func SET_WITH_Expiry(key string, value string, expiry time.Time) {
	store[key] = &StoreItem{
		Value:  value,
		Expiry: expiry,
	}
}

func PerformPong(args []string) string {
	return responses.StringMsg("PONG")
}

func PerformSet(args []string) string {
	if len(args) < 2 {
		return responses.ErrorMsg("invalid syntax provided to 'SET'")
	}

	var exp time.Time
	key, val := &args[0], &args[1]

	if len(args) > 2 {
		position := 2
		for position < len(args) {
			switch strings.ToLower(args[position]) {
			// milliseconds
			case "px":
				if len(args) < position+1 {
					return responses.ErrorMsg("no time provided to 'PX'")
				}
				expMillis, err := strconv.Atoi(string(args[position+1]))
				if err != nil {
					return responses.ErrorMsg("invalid format provided to 'PX'")
				}

				exp = time.Now().Add(time.Millisecond * time.Duration(expMillis))
				position += 2

			// seconds
			case "ex":
				if len(args) < position+1 {
					return responses.ErrorMsg("no time provided to 'EX'")
				}
				expSeconds, err := strconv.Atoi(string(args[position+1]))
				if err != nil {
					return responses.ErrorMsg("invalid format provided to 'EX'")
				}

				exp = time.Now().Add(time.Second * time.Duration(expSeconds))
				position += 2
			default:
				return responses.ErrorMsg(fmt.Sprintf("invalid argument '%s'", args[position]))
			}
		}
	}
	// If no expiry is set, set the value without expiry
	if exp == (time.Time{}) {
		SET(*key, *val)
	} else {
		SET_WITH_Expiry(*key, *val, time.Now().Add(time.Second*time.Duration(exp.Unix())))
	}
	return responses.StringMsg("OK")
}

// PerformGet retrieves a value from the database,
// if it exists and is not expired. If it is expired, it will be deleted
func PerformGet(args []string) string {
	if len(args) < 1 {
		return responses.ErrorMsg("no value provided to 'GET'")
	}

	item, exists := store[args[0]]

	if !exists || item == nil {
		return responses.ErrorMsg(fmt.Sprintf("no value found for key '%s'", args[0]))
	}

	// Enforce mutual exclusion on the expiry operation and retrieval
	item.Mutex.Lock()
	defer item.Mutex.Unlock()

	now := time.Now()
	if !item.Expiry.IsZero() && item.Expiry.Before(now) {
		delete(store, args[0])
		return responses.ErrorMsg(fmt.Sprintf("no value found for key '%s'", args[0]))
	}

	return responses.StringMsg(item.Value)
}

// PerformDel deletes a value from the database
// if it exists. If it does not exist, it returns an error message
func PerformDel(args []string) string {
	if len(args) < 1 {
		return responses.ErrorMsg("no value provided to 'DEL'")
	}

	item, exists := store[args[0]]
	if exists {
		item.Mutex.Lock()
		defer item.Mutex.Unlock()
		delete(store, args[0])
		return responses.StringMsg("The key has been deleted")
	}

	return responses.ErrorMsg(fmt.Sprintf("no value found for key '%s'", args[0]))
}

func PerformExists(args []string) string {
	if len(args) < 1 {
		return responses.ErrorMsg("no value provided to 'EXISTS'")
	}
	answer := PerformGet(args)
	if strings.Contains(answer, "no value found") {
		return responses.StringMsg("False")
	}
	return responses.StringMsg("True")
}

// PerformTTL returns the remaining time to live for a key in seconds
func PerformTTL(args []string) string {
	if len(args) != 1 {
		return responses.ErrorMsg("wrong number of arguments for 'TTL' command")
	}

	key := args[0]
	remaining, exists := store[key]
	if !exists {
		return responses.ErrorMsg(fmt.Sprintf("no value found for key '%s'", key))
	}

	// Enforce mutual exclusion on the expiry operation and retrieval
	remaining.Mutex.Lock()
	defer remaining.Mutex.Unlock()

	// If the key has no expiry, return 0
	if remaining.Expiry.IsZero() {
		return responses.StringMsg("0")
	}

	now := time.Now()
	// If the key is expired, delete it and return -2
	if remaining.Expiry.Before(now) {
		delete(store, key)
		return responses.StringMsg("-2")
	}

	ttl := remaining.Expiry.Sub(now).Seconds()
	return responses.StringMsg(fmt.Sprintf("%d", int(ttl)))
}

// PerformGETDEL retrieves a value and deletes it in a single operation
func PerformGETDEL(args []string) string {
	if len(args) != 1 {
		return responses.ErrorMsg("wrong number of arguments for 'GETDEL' command")
	}

	key := args[0]
	item, exists := store[key]
	if !exists {
		return responses.ErrorMsg(fmt.Sprintf("no value found for key '%s'", key))
	}

	// Enforce mutual exclusion on the expiry operation and retrieval
	item.Mutex.Lock()
	defer item.Mutex.Unlock()

	now := time.Now()
	// If the key is expired, delete it and return -2
	if !item.Expiry.IsZero() && item.Expiry.Before(now) {
		delete(store, key)
		return responses.StringMsg("-2")
	}

	value := item.Value
	delete(store, key)
	return responses.StringMsg(value)
}

// PerformRename renames a key to a new key
func PerformRename(args []string) string {
	if len(args) != 2 {
		return responses.ErrorMsg("wrong number of arguments for 'RENAME' command")
	}

	oldKey := args[0]
	newKey := args[1]

	item, exists := store[oldKey]
	if !exists {
		return responses.ErrorMsg(fmt.Sprintf("no value found for key '%s'", oldKey))
	}

	// Enforce mutual exclusion on the expiry operation and retrieval
	item.Mutex.Lock()
	defer item.Mutex.Unlock()

	store[newKey] = item
	delete(store, oldKey)

	return responses.StringMsg("OK")
}

// PerformExpire sets an expiry time for a key
func PerformExpire(args []string) string {
	if len(args) != 2 {
		return responses.ErrorMsg("wrong number of arguments for 'EXPIRE' command")
	}

	key := args[0]
	expirySeconds, err := strconv.Atoi(args[1])
	if err != nil {
		return responses.ErrorMsg("invalid expiry time provided")
	}

	item, exists := store[key]
	if !exists {
		return responses.ErrorMsg(fmt.Sprintf("no value found for key '%s'", key))
	}

	// Enforce mutual exclusion on the expiry operation and retrieval
	item.Mutex.Lock()
	defer item.Mutex.Unlock()

	item.Expiry = time.Now().Add(time.Second * time.Duration(expirySeconds))
	return responses.StringMsg("OK")
}

func WatchCommands(args []string) string {
	if len(args) != 0 {
		return responses.ErrorMsg("no arguments expected for 'WATCH' command")
	}
	url := "watchCommandsPageUrl"
	return responses.StringMsg("Fetching watch commands from: " + url)
}
