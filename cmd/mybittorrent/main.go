package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
	// bencode "github.com/jackpal/bencode-go" // Available if you need it!
)

// Ensures gofmt doesn't remove the "os" encoding/json import (feel free to remove this!)
var _ = json.Marshal

// Example:
// - 5:hello -> hello
// - 10:hello12345 -> hello12345
func decodeBencode(bencodedString string) (interface{}, error) {
	// decode strings
	if unicode.IsDigit(rune(bencodedString[0])) {
		var firstColonIndex int

		for i := 0; i < len(bencodedString); i++ {
			if bencodedString[i] == ':' {
				firstColonIndex = i
				break
			}
		}

		lengthStr := bencodedString[:firstColonIndex]

		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return "", err
		}

		return bencodedString[firstColonIndex+1 : firstColonIndex+1+length], nil

		// Decode integers
	} else if string(bencodedString[0]) == "i" {
		var decimalNumber int
		var err error

		number := bencodedString[1:len(bencodedString)-1]

		if string(number[0]) == "-" {
			decimalNumber, err = strconv.Atoi(number[1:])
			if err != nil {
				return "", err
			}
			decimalNumber *= -1
		} else {
			decimalNumber, err =  strconv.Atoi(number)
			if err != nil {
				return "", err
			}
		}

		return decimalNumber, nil

		// Decode lists
	}  else if bencodedString[0] == 'l'{
		encodedList := bencodedString[1:len(bencodedString)-1]
		var err error
		var item interface{}
		var list []interface{}

		for len(encodedList) > 0 {
			item, err = decodeBencode(encodedList)
			if err != nil {
				return "", err
			}
			list = append(list, item)
			_, encodedList, _ = splitEncodedItem(encodedList)
		}
		
		if list == nil {
			return []string{}, nil
		}
		return list, nil
	} else {
		return "", fmt.Errorf("only strings are supported at the moment")
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	command := os.Args[1]

	if command == "decode" {
		// Uncomment this block to pass the first stage
		
		bencodedValue := os.Args[2]
		
		decoded, err := decodeBencode(bencodedValue)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}

func splitEncodedItem(s string) (string, string, error) {
	if len(s) == 0 {
		return "", "", fmt.Errorf("empty string")
	}

	switch s[0] {
	case 'i':
		endIndex := strings.IndexByte(s, 'e')
		
		if endIndex == -1 {
			return "", "", fmt.Errorf("malformed integer")
		}
		return s[:endIndex+1], s[endIndex+1:], nil

	case 'l':
		depth := 0

		for i, char := range s {
			if char == 'l' {
				depth++
			} else if char == 'e' {
				depth--
				if depth == 0 {
					return s[:i+1], s[i+1:], nil
				}
			}
		}
		return "", "", fmt.Errorf("malformed list")

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		colonIndex := strings.IndexByte(s, ':')
		if colonIndex == -1 {
			return "", "", fmt.Errorf("malformed string")
		}
		length, err := strconv.Atoi(s[:colonIndex])
		if err!= nil {
			return "", "", err
		}
		endIndex := colonIndex + length + 1
		if endIndex > len(s) {
			return "", "", fmt.Errorf("string length exceeds available data")
		}
		return s[:endIndex], s[endIndex:], nil

	default:
		return "", "", fmt.Errorf("unsupported type identifier : %c", s[0])
	}
}