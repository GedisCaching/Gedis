package RESP

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type StoreItem struct {
	Value  interface{}
	Expiry time.Time
	Mutex  sync.Mutex
}

// Example usage of the store to avoid unused variable error
var store = make(map[string]*StoreItem)

func SET(key string, value interface{}) {
	store[key] = &StoreItem{
		Value: value,
	}
}

func SET_WITH_Expiry(key string, value interface{}, expiry time.Time) {
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

	if exp == (time.Time{}) {
		SET(*key, *val)
	} else {
		SET_WITH_Expiry(*key, *val, exp)
	}
	return stringMsg("OK")
}
