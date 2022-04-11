package ommerrors

import (
	"fmt"
	"strconv"
)

type ReplaceConflictError struct{ Text string }

func NewReplaceConflictError(file string, line int) *ReplaceConflictError {
	return &ReplaceConflictError{
		Text: fmt.Sprintf("ReplaceConflictError: file=%v , line=%v", file, strconv.Itoa(line)),
	}
}

func (e *ReplaceConflictError) Error() string {
	return e.Text
}
