package validator

import (
	"errors"
	"fmt"
	"strings"
)

func ValidateCurrencyCode(code string) error {
	if len(code) == 0 {
		return errors.New("Currency code is not present in the request")
	}
	if len(code) != 3 {
		return errors.New(fmt.Sprintf(
			"Currency code must contain exactly 3 letters as defined in ISO 4217, got: %s", code,
		))
	}
	if code != strings.ToUpper(code) {
		return errors.New(fmt.Sprintf(
			"Currency code must contain exactly 3 uppercase letters as defined in ISO 4217, got: %s", code,
		))
	}
	return nil
}
