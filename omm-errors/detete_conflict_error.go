package ommerrors

import (
	"fmt"
	"strconv"
)

type DeleteConflictError struct{ Text string }

func NewDeleteConflictError(file string, line int) *DeleteConflictError {
	return &DeleteConflictError{
		Text: fmt.Sprintf("DeleteConflictError: file=%v , line=%v", file, strconv.Itoa(line)),
	}
}

func (e *DeleteConflictError) Error() string {
	return e.Text
}
