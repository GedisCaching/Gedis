package responses

import (
	"fmt"
)

func StringMsg(msg string) string {
	return "+" + msg
}

func ErrorMsg(msg string) string {
	return fmt.Sprintf("-ERR %s", msg)
}

func NilBulkStringMsg() string {
    return "$-1"
}