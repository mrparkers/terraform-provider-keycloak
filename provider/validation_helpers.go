package provider

import (
	"fmt"
)

func logicKeyValidation(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	if v != "POSITIVE" && v != "NEGATIVE" {
		errs = append(errs, fmt.Errorf("%q must 'POSITIVE' or 'NEGATIVE' %d", key, val))
	}
	return
}
