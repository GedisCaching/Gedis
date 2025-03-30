package RESP

import (
	"fmt"
	"strings"
)

type Parser struct {
	NumberOfExpectedArguments int
	LengthOfNextArgument      int
	NumberOfArgumentsParsed   int
}

func Parse(command []byte) string {
	var cmd string
	args := []string{}

	parser := &Parser{
		NumberOfExpectedArguments: 0,
		LengthOfNextArgument:      0,
		NumberOfArgumentsParsed:   0,
	}

	position := 0
	for position < len(command) {
		switch command[position] {
		// Array of Strings
		case '*':
			result, err := getIntArg(position+1, command)
			if err != nil {
				return err.Error()
			}
			parser.NumberOfExpectedArguments = result.Result
			args = make([]string, parser.NumberOfExpectedArguments-1)
			position += result.PositionsParsed

		// Bulk String
		case '$':
			result, err := getIntArg(position+1, command)
			if err != nil {
				return err.Error()
			}
			// Enforce the length of the next argument
			parser.LengthOfNextArgument = result.Result
			position += result.PositionsParsed

		// to denote the end of the argument count
		case '\r':
			position++
		// to denote the end of the argument count
		case '\n':
			position++

		default:
			if parser.NumberOfExpectedArguments == 0 || parser.LengthOfNextArgument == 0 || parser.NumberOfArgumentsParsed >= parser.NumberOfExpectedArguments {
				return stringMsg("Invalid syntax")
			}

			// parse it!
			parsedItem := string(command[position : parser.LengthOfNextArgument+position])
			if parser.NumberOfArgumentsParsed > 0 {
				// it's an arg, add to args array
				args[parser.NumberOfArgumentsParsed-1] = parsedItem
			} else {
				// The first 'arg' we parse is the primary command
				cmd = parsedItem
			}

			position += parser.LengthOfNextArgument

			parser.LengthOfNextArgument = 0
			parser.NumberOfArgumentsParsed += 1
		}
	}
	return ParseCommand(cmd, args)
}

func ParseCommand(command string, args []string) string {
	cmd := strings.ToUpper(command)
	fmt.Printf("Received '%s' command\n", cmd)

	switch cmd {
	case "PING":
		return PerformPong(args)
	case "SET":
		return PerformSet(args)
	case "GET":
		return PerformGet(args)
	case "DEL":
		return PerformDel(args)
	case "EXISTS":
		return PerformExists(args)
	default:
		return errorMsg(fmt.Sprintf("unknown command '%s'", cmd))
	}
}
