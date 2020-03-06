package provider

import (
	"fmt"
)

func logicKeyValidation(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	errMsg := fmt.Errorf("%q must 'POSITIVE' or 'NEGATIVE' %d", key, val)
	if v != "POSITIVE" {
		errs = append(errs, errMsg)
	}
	if v != "NEGATIVE" {
		errs = append(errs, errMsg)
	}
	return
}
