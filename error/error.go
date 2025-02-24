package error

import (
	"context"
	"errors"
	"fmt"
)

// HandleDeadlineExceededError checks if the given error is a context deadline exceeded error.
func HandleDeadlineExceededError(err error) error {
	if err == context.DeadlineExceeded {
		message := fmt.Sprintf("Operation timed out: %v", err)
		return errors.New(message)
	}
	return nil
}
