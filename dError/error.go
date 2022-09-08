package dError

import "fmt"

type errorType struct {
	PanicCode string
	PanicMsg  error
}

func NewError(PanicCode string, panicMsg error) *errorType {
	return &errorType{
		PanicCode: PanicCode,
		PanicMsg:  panicMsg,
	}
}

func (e *errorType) Error() string {
	return fmt.Sprintf("(%s) %s", e.PanicCode, e.PanicMsg.Error())
}
