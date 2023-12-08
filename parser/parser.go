package parser

import (
	"encoding/json"
	"io"
)

func FromJSON[T any](r io.Reader) T {
	var parsedData T
	json.NewDecoder(r).Decode(&parsedData)
	return parsedData
}
