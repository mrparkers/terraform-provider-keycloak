package keycloak

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type KeycloakBoolQuoted bool

func (c KeycloakBoolQuoted) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if c == false {
		buf.WriteString(`""`)
	} else {
		buf.WriteString(strconv.Quote(strconv.FormatBool(bool(c))))
	}
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

func GetTags(f interface{}) []string {
	var result []string = make([]string, 0)

	t := reflect.ValueOf(f).Elem().Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := string(field.Tag.Get("json"))
		if tag != "" {
			result = append(result, tag)
		}
	}

	return result
}

func GetReservedKeys(f interface{}) map[string]bool {
	var result = map[string]bool{}

	tags := GetTags(f)
	for i := 0; i < len(tags); i++ {
		result[tags[i]] = true
	}

	return result
}

func getIdFromLocationHeader(locationHeader string) string {
	parts := strings.Split(locationHeader, "/")

	return parts[len(parts)-1]
}

// Converts duration string to a string representing the number of milliseconds, which is used by the Keycloak API
// Ex: "1h" => "3600000"
func getMillisecondsFromDurationString(s string) (string, error) {
	duration, err := time.ParseDuration(s)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(duration.Seconds() * 1000)), nil
}

// Converts a string representing milliseconds from Keycloak API to a duration string used by the provider
// Ex: "3600000" => "1h0m0s"
func GetDurationStringFromMilliseconds(milliseconds string) (string, error) {
	ms, err := strconv.Atoi(milliseconds)
	if err != nil {
		return "", err
	}

	return (time.Duration(ms) * time.Millisecond).String(), nil
}

func parseBoolAndTreatEmptyStringAsFalse(b string) (bool, error) {
	if b == "" {
		return false, nil
	}

	return strconv.ParseBool(b)
}
