package RESP

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
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
	return stringMsg("PONG")
}

func PerformSet(args []string) string {
	if len(args) < 2 {
		return errorMsg("invalid syntax provided to 'SET'")
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
					return errorMsg("no time provided to 'PX'")
				}
				expMillis, err := strconv.Atoi(string(args[position+1]))
				if err != nil {
					return errorMsg("invalid format provided to 'PX'")
				}

				exp = time.Now().Add(time.Millisecond * time.Duration(expMillis))
				position += 2

			// seconds
			case "ex":
				if len(args) < position+1 {
					return errorMsg("no time provided to 'EX'")
				}
				expSeconds, err := strconv.Atoi(string(args[position+1]))
				if err != nil {
					return errorMsg("invalid format provided to 'EX'")
				}

				exp = time.Now().Add(time.Second * time.Duration(expSeconds))
				position += 2
			default:
				return errorMsg(fmt.Sprintf("invalid argument '%s'", args[position]))
			}
		}
	}
	// If no expiry is set, set the value without expiry
	if exp == (time.Time{}) {
		SET(*key, *val)
	} else {
		SET_WITH_Expiry(*key, *val, exp)
	}
	return stringMsg("OK")
}

// PerformGet retrieves a value from the database,
// if it exists and is not expired. If it is expired, it will be deleted
func PerformGet(args []string) string {
	if len(args) < 1 {
		return errorMsg("no value provided to 'GET'")
	}

	item, exists := store[args[0]]

	if !exists || item == nil {
		return errorMsg(fmt.Sprintf("no value found for key '%s'", args[0]))
	}

	// Enforce mutual exclusion on the expiry operation and retrieval
	item.Mutex.Lock()
	defer item.Mutex.Unlock()

	now := time.Now()
	if !item.Expiry.IsZero() && item.Expiry.Before(now) {
		delete(store, args[0])
		return errorMsg(fmt.Sprintf("no value found for key '%s'", args[0]))
	}

	return stringMsg(item.Value)
}

// PerformDel deletes a value from the database
// if it exists. If it does not exist, it returns an error message
func PerformDel(args []string) string {
	if len(args) < 1 {
		return errorMsg("no value provided to 'DEL'")
	}

	item, exists := store[args[0]]
	if exists {
		item.Mutex.Lock()
		defer item.Mutex.Unlock()
		delete(store, args[0])
		return stringMsg("The key has been deleted")
	}

	return errorMsg(fmt.Sprintf("no value found for key '%s'", args[0]))
}

func PerformExists(args []string) string {
	if len(args) < 1 {
		return errorMsg("no value provided to 'EXISTS'")
	}
	answer := PerformGet(args)
	if strings.Contains(answer, "no value found") {
		return stringMsg("False")
	}
	return stringMsg("True")
}
