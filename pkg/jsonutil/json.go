package jsonutil

import (
	json "github.com/bytedance/sonic"
)

// ParseJSON parses a JSON string into the result.
func ParseJSON(data string, result any) error {
	return ParseJSONFromBytes([]byte(data), result)
}

// StringifyJSON converts data to a JSON string.
func StringifyJSON(data any) (string, error) {
	b, err := StringifyJSONToBytes(data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ParseJSONFromBytes parses JSON bytes into the result.
func ParseJSONFromBytes(data []byte, result any) error {
	return json.Unmarshal(data, result)
}

// StringifyJSONToBytes converts data to JSON bytes.
func StringifyJSONToBytes(data any) ([]byte, error) {
	return json.Marshal(&data)
}
