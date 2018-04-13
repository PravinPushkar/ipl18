package errors

import "fmt"

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type DaoError struct {
	Code      int
	ActualErr error
	UserErr   error
}

func (d *DaoError) Error() string {
	return fmt.Sprintf("Code - %d,Err - %s", d.Code, map[bool]string{false: d.ActualErr.Error(), true: d.UserErr.Error()}[d.UserErr == nil])
}

func (e *Error) Error() string {
	return fmt.Sprintf("Code - %d,Err - %s", e.Code, e.Message)
}
