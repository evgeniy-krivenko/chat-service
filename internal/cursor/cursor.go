package cursor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func Encode(data any) (string, error) {
	d, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("encode cursor: %v", err)
	}
	return base64.StdEncoding.EncodeToString(d), nil
}

func Decode(in string, to any) error {
	decoded, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return fmt.Errorf("decode cursor: %v", err)
	}

	return json.Unmarshal(decoded, to)
}
