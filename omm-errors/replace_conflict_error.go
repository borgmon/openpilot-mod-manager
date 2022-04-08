package ommerrors

import "strconv"

type ReplaceConflictError struct{ Text string }

func NewReplaceConflictError(file string, line int) *ReplaceConflictError {
	return &ReplaceConflictError{
		Text: "ReplaceConflictError: file=" + file + " , line=" + strconv.Itoa(line),
	}
}

func (e *ReplaceConflictError) Error() string {
	return e.Text
}
