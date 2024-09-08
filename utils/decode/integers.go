package decode

import (
	"fmt"
	"strconv"
)

// ToInteger decodes an integer
// eg. i123e -> 123
func ToInteger(data string) (int64, error) {
	if data[0] != 'i' || data[len(data)-1] != 'e' {
		return 0, fmt.Errorf("invalid integer")
	}

	return strconv.ParseInt(data[1:len(data)-1], 10, 64)
}
