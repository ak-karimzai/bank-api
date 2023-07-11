package validators

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[A-Za-z ]+$`).MatchString
)

func ValidateUsername(val string) error {
	if err := ValidateString(val, 3, 100); err != nil {
		return err
	}

	if !isValidUsername(val) {
		return fmt.Errorf("must contains only letters, digits, or underscore")
	}
	return nil
}

func ValidatePwd(val string) error {
	return ValidateString(val, 6, 200)
}

func ValidateEmail(val string) error {
	if err := ValidateString(val, 3, 200); err != nil {
		return err
	}
	_, err := mail.ParseAddress(val)
	if err != nil {
		return fmt.Errorf("%v is not a valid email address", val)
	}
	return nil
}

func ValidateFullName(val string) error {
	if err := ValidateString(val, 3, 100); err != nil {
		return err
	}

	if !isValidFullName(val) {
		return fmt.Errorf("must contain only letters or spaces")
	}
	return nil
}
