package decode

import (
	"fmt"
	"strconv"
	"strings"
)

func ToList(data string) ([]string, error) {
	if len(data) < 2 || data[0] != 'l' || data[len(data)-1] != 'e' {
		return nil, fmt.Errorf("invalid list format")
	}

	data = data[1 : len(data)-1]
	if len(data) == 0 {
		return []string{}, nil
	}

	var result []string

	for len(data) > 0 {
		colon := strings.Index(data, ":")
		if colon == -1 {
			return nil, fmt.Errorf("invalid list: missing colon")
		}

		length, err := strconv.Atoi(data[:colon])
		if err != nil || length < 0 {
			return nil, fmt.Errorf("invalid length in list")
		}

		if len(data) < colon+1+length {
			return nil, fmt.Errorf("invalid list: data too short")
		}

		item := data[colon+1 : colon+1+length]
		result = append(result, item)

		data = data[colon+1+length:]
	}

	return result, nil
}
