package tmerr

import "fmt"

type TimeManError string

var (
	InputParamsError         TimeManError = "InputParamsError"
	NullStringParamsError    error        = fmt.Errorf("NullStringParamsError(%w)", InputParamsError)
	TooLongStringParamsError error        = fmt.Errorf("TooLongStringParamsError(%w)", InputParamsError)
)

func Btw(err error, extra string) error {
	return fmt.Errorf("%s: %w", extra, err)

}
