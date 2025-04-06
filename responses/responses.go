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

// ArrayMsg formats a slice of strings as a RESP array
func ArrayMsg(elements []string) string {
	result := fmt.Sprintf("*%d\r\n", len(elements))
	for _, element := range elements {
		result += fmt.Sprintf("$%d\r\n%s\r\n", len(element), element)
	}
	return result
}

// IntegerMsg formats an integer as a RESP integer
func IntegerMsg(n int) string {
	return fmt.Sprintf(":%d", n)
}
