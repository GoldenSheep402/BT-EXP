package decode

import (
	"fmt"
	"strconv"
)

// ToString decodes a string
func ToString(data string) (string, error) {
	length := 0
	for i := 0; i < len(data); i++ {
		if data[i] == ':' {
			length, _ = strconv.Atoi(data[:i])
			break
		}
	}

	if length == 0 {
		return "", fmt.Errorf("invalid string")
	}

	return data[len(data)-length:], nil
}
