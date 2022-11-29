package types

import (
	"bytes"
	"strconv"
)

// KeycloakBoolQuoted is a bool that is marshalled to a quoted string in JSON. This is needed for some boolean
// attributes in the Keycloak API that are treated as strings.
type KeycloakBoolQuoted bool

func (c KeycloakBoolQuoted) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(strconv.Quote(strconv.FormatBool(bool(c))))
	return buf.Bytes(), nil
}

func (c *KeycloakBoolQuoted) UnmarshalJSON(in []byte) error {
	value := string(in)
	if value == `""` {
		*c = false
		return nil
	}
	unquoted, err := strconv.Unquote(value)
	if err != nil {
		return err
	}
	var b bool
	b, err = strconv.ParseBool(unquoted)
	if err != nil {
		return err
	}
	res := KeycloakBoolQuoted(b)
	*c = res
	return nil
}
