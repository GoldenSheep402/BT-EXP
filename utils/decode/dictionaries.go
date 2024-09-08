package decode

import (
	"fmt"
	"strconv"
	"strings"
)

// ToDictionaries decodes a bencoded dictionary
// Example 1: d3:cow3:moo4:spam4:eggse => {'cow': 'moo', 'spam': 'eggs'}
// Example 2: d4:spaml1:a1:bee => {'spam': ['a', 'b']}
func ToDictionaries(data string) (map[string]interface{}, error) {
	if data[0] != 'd' || data[len(data)-1] != 'e' {
		return nil, fmt.Errorf("invalid dictionary")
	}

	data = data[1 : len(data)-1]
	result := make(map[string]interface{})

	for len(data) > 0 {
		// Find the first colon which separates the key length
		colon := strings.Index(data, ":")
		if colon == -1 {
			return nil, fmt.Errorf("invalid dictionary: missing colon")
		}

		// Get the key length
		keyLength, err := strconv.Atoi(data[:colon])
		if err != nil || keyLength < 0 {
			return nil, fmt.Errorf("invalid key length in dictionary")
		}

		// Get the key
		if len(data) < colon+1+keyLength {
			return nil, fmt.Errorf("invalid dictionary: data too short for key")
		}
		key := data[colon+1 : colon+1+keyLength]
		data = data[colon+1+keyLength:]

		if len(data) == 0 {
			result[key] = ""
			return result, nil
		}

		switch data[0] {
		case 'l':
			// Handle list
			endList := strings.Index(data, "e")
			if endList == -1 {
				return nil, fmt.Errorf("invalid list in dictionary")
			}
			list, err := ToList(data[:endList+1])
			if err != nil {
				return nil, err
			}
			result[key] = list
			data = data[endList+1:]
		case 'd':
			endDict := strings.Index(data, "e")
			if endDict == -1 {
				return nil, fmt.Errorf("invalid nested dictionary")
			}
			dict, err := ToDictionaries(data[:endDict+1])
			if err != nil {
				return nil, err
			}
			result[key] = dict
			data = data[endDict+1:]
		default:
			colon = strings.Index(data, ":")
			if colon == -1 {
				return nil, fmt.Errorf("invalid string in dictionary")
			}
			valueLength, err := strconv.Atoi(data[:colon])
			if err != nil || valueLength < 0 {
				return nil, fmt.Errorf("invalid value length in dictionary")
			}
			if len(data) < colon+1+valueLength {
				return nil, fmt.Errorf("invalid dictionary: data too short for value")
			}
			value := data[colon+1 : colon+1+valueLength]
			result[key] = value
			data = data[colon+1+valueLength:]
		}
	}

	return result, nil
}
