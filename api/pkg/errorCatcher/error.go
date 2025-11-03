package errorCatcher

import (
	"fmt"
)

// PanicIfErr func
func PanicIfErr(err, errorType, errorSubject error) {
	if err != nil {
		panic(ConcatError(errorType, errorSubject, err))
	}
}

// ReturnIfErr func
func ReturnIfErr(err, errorType, errorSubject error) error {
	if err != nil {
		return ConcatError(errorType, errorSubject, err)
	}
	return nil
}

// ConcatError func
func ConcatError(errorType, errorSubject, err error) error {
	return fmt.Errorf("%w: %w: %w", errorType, errorSubject, err)
}
