package e

import "fmt"

func WrapError(errMsg string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", errMsg, err)
}