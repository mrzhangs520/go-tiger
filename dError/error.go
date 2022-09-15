package dError

import "errors"

type errorType struct {
	PanicCode    string
	SourceErr    error
	PanicUserMag string
}

func NewError(panicUserMag, panicCode string, sourceErrList ...error) *errorType {
	var err error
	if 0 == len(sourceErrList) {
		err = errors.New("")
	} else {
		err = sourceErrList[0]
	}
	return &errorType{
		PanicUserMag: panicUserMag,
		PanicCode:    panicCode,
		SourceErr:    err,
	}
}

func (e *errorType) Error() string {
	return e.PanicUserMag
}

func (e *errorType) GetContent() *errorType {
	return e
}
