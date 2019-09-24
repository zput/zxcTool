package zxcUtil

import (
	"encoding/json"
	"fmt"
)

// interface
type ErrorInterface interface {
	Status() int
}

//struct
type Error struct {
	status  int    `json:"-"`
	Code    *int   `json:"code,omitempty"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

func (e Error) Status() int {
	return e.status
}

//new a above error
func CodedErrorf(status int, code *int, format string, args ...interface{}) error {
	return &Error{status: status, Code: code, Message: fmt.Sprintf(format, args...)}
}

func Errorf(status int, format string, args ...interface{}) error {
	return &Error{status: status, Message: fmt.Sprintf(format, args...)}
}
