package RESP

import (
	"errors"
	"strconv"
	"strings"
)

type GetIntArgResult struct {
	Result          int // parsing operation
	PositionsParsed int // How many indexes were parsed
}

// getIntArg parses an integer from the byte array starting at the given position
// and returns the result and the number of positions parsed.
func getIntArg(startPosition int, arr []byte) (*GetIntArgResult, error) {
	var err error
	result := &GetIntArgResult{
		Result:          0,
		PositionsParsed: 0,
	}

	// Get digits until termination characters
	notDone := true
	position := startPosition

	var stringVal strings.Builder
	for position < len(arr) && notDone {
		// Check for literal "\r\n" string (4 characters)
		if position+3 < len(arr) && arr[position] == '\\' && arr[position+1] == 'r' && arr[position+2] == '\\' && arr[position+3] == 'n' {
			notDone = false
			result.PositionsParsed += 4 // Skip over "\r\n"
		} else if position+1 < len(arr) && arr[position] == '\r' && arr[position+1] == '\n' {
			notDone = false
			result.PositionsParsed += 2 // Skip over CR+LF
		} else {
			stringVal.WriteByte(arr[position])
			result.PositionsParsed += 1
			position += 1
		}
	}

	if stringVal.Len() == 0 {
		return nil, errors.New("no value was detected")
	}

	resultInt, err := strconv.Atoi(stringVal.String())
	if err != nil {
		return nil, errors.New("failed to parse int")
	}

	result.Result = resultInt
	return result, nil
}
