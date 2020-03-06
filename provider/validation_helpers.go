package provider

import (
	"fmt"
)

func logicKeyValidation(val interface{}, key string) (warns []string, errs []error) {
	v := val.(string)
	isNotPositive := v != "POSITIVE"
	isNotNegative := v != "NEGATIVE"
	if isNotPositive || isNotNegative {
		errs = append(errs, fmt.Errorf("%q must 'POSITIVE' or 'NEGATIVE' %d", key, val))
	}
	return
}
