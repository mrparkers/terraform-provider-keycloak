package types

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// KeycloakSliceQuoted is a slice of strings that is marshaled to a quoted JSON string array
type KeycloakSliceQuoted []string

func (s KeycloakSliceQuoted) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if s == nil || len(s) == 0 {
		buf.WriteString(`""`)
	} else {
		sliceAsString := make([]string, len(s))
		for i, v := range s {
			sliceAsString[i] = v
		}

		stringAsJSON, err := json.Marshal(sliceAsString)
		if err != nil {
			return nil, err
		}

		buf.WriteString(strconv.Quote(string(stringAsJSON)))
	}

	return buf.Bytes(), nil
}
