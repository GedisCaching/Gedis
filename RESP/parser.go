package RESP

import (
	"fmt"
	"strings"
	responses "github.com/GedisCaching/Gedis/responses"
)

type Parser struct {
	NumberOfExpectedArguments int
	LengthOfNextArgument      int
	NumberOfArgumentsParsed   int
}

func Parse(command []byte) string {
	// Check if this is a plain text command (doesn't start with RESP markers)
	if len(command) > 0 && command[0] != '*' && command[0] != '$' {
		return parsePlainTextCommand(command)
	}

	// Original RESP parsing logic
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
				return responses.ErrorMsg(err.Error())
			}
			parser.NumberOfExpectedArguments = result.Result
			args = make([]string, parser.NumberOfExpectedArguments-1)
			position += result.PositionsParsed

		// Bulk String
		case '$':
			result, err := getIntArg(position+1, command)
			if err != nil {
				return responses.ErrorMsg(err.Error())
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
				return responses.StringMsg("Invalid syntax")
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

// parsePlainTextCommand handles commands in plain text format like "SET key value EX 30"
func parsePlainTextCommand(command []byte) string {
	// Trim any trailing whitespace, CR, LF
	commandStr := strings.TrimSpace(string(command))

	// Split the command by whitespace
	parts := strings.Fields(commandStr)
	if len(parts) == 0 {
		return responses.ErrorMsg("empty command")
	}

	// The first part is the command, the rest are arguments
	cmd := parts[0]
	args := []string{}
	if len(parts) > 1 {
		args = parts[1:]
	}

	fmt.Printf("Parsed plain text command: %s, args: %v\n", cmd, args)

	// Pass to the command handler
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
		return responses.ErrorMsg(fmt.Sprintf("unknown command '%s'", cmd))
	}
}
