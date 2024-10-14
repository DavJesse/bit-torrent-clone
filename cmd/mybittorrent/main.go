package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
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

		}

		func splitEncodedItem(s string) (string, string, err) {

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
