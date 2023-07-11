package validators

import (
	"fmt"
)

func ValidateString(val string, min, max int) error {
	n := len(val)
	if n < min || n > max {
		return fmt.Errorf("must contain %d-%d charecters", min, max)
	}
	return nil
}
