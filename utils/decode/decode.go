package decode

import (
	"errors"
	"reflect"
	"strconv"
)

func decode(data []byte, value interface{}) error {
	target := reflect.ValueOf(value)
	if target.Kind() != reflect.Ptr || target.IsNil() {
		return errors.New("[invalid target value]")
	}

	result, rest, err := decodeValue(data)
	if err != nil {
		return err
	}
	if len(rest) > 0 {
		return errors.New("[unexpected data after end of Bencode]")
	}

	target.Elem().Set(reflect.ValueOf(result))
	return nil
}

func decodeString(data []byte) (string, []byte, error) {
	colonPos := 0
	for colonPos < len(data) && data[colonPos] != ':' {
		colonPos++
	}
	if colonPos == len(data) {
		return "", nil, errors.New("[invalid string format]")
	}
	length, err := strconv.Atoi(string(data[:colonPos]))
	if err != nil {
		return "", nil, err
	}
	if len(data)-colonPos-1 < length {
		return "", nil, errors.New("[invalid string format]")
	}
	return string(data[colonPos+1 : colonPos+1+length]), data[colonPos+1+length:], nil
}

func decodeInt(data []byte) (int64, []byte, error) {
	if len(data) == 0 || data[0] != 'i' {
		return 0, nil, errors.New("[invalid integer format]")
	}
	data = data[1:]
	endPos := 0
	for endPos < len(data) && data[endPos] != 'e' {
		endPos++
	}
	if endPos == len(data) {
		return 0, nil, errors.New("[invalid integer format]")
	}
	value, err := strconv.ParseInt(string(data[:endPos]), 10, 64)
	if err != nil {
		return 0, nil, err
	}
	return value, data[endPos+1:], nil
}

func decodeList(data []byte) ([]interface{}, []byte, error) {
	if len(data) == 0 || data[0] != 'l' {
		return nil, nil, errors.New("[invalid list format]")
	}
	data = data[1:]
	result := make([]interface{}, 0)
	for len(data) > 0 && data[0] != 'e' {
		item, rest, err := decodeValue(data)
		if err != nil {
			return nil, nil, err
		}
		result = append(result, item)
		data = rest
	}
	if len(data) == 0 || data[0] != 'e' {
		return nil, nil, errors.New("[invalid list format]")
	}
	return result, data[1:], nil
}

func decodeDict(data []byte) (map[string]interface{}, []byte, error) {
	if len(data) == 0 || data[0] != 'd' {
		return nil, nil, errors.New("[invalid dict format]")
	}
	data = data[1:]
	result := make(map[string]interface{})
	for len(data) > 0 && data[0] != 'e' {
		key, rest, err := decodeString(data)
		if err != nil {
			return nil, nil, err
		}
		value, rest, err := decodeValue(rest)
		if err != nil {
			return nil, nil, err
		}
		result[key] = value
		data = rest
	}
	if len(data) == 0 || data[0] != 'e' {
		return nil, nil, errors.New("[invalid dict format]")
	}
	return result, data[1:], nil
}

func decodeValue(data []byte) (interface{}, []byte, error) {
	if len(data) == 0 {
		return nil, nil, errors.New("[invalid data format]")
	}
	switch data[0] {
	case 'i':
		return decodeInt(data)
	case 'l':
		return decodeList(data)
	case 'd':
		return decodeDict(data)
	default:
		return decodeString(data)
	}
}
