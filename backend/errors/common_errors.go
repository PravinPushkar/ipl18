package errors

import "fmt"

var (
	ErrDBIssue           = fmt.Errorf("could not connect to db")
	ErrParseContext      = fmt.Errorf("could not parse context")
	ErrTokenInfoMismatch = fmt.Errorf("token info mismatch")
	ErrParseRequest      = fmt.Errorf("could not parse request")
	ErrEncodingResponse  = fmt.Errorf("error encoding response")
	ErrGettingToken      = fmt.Errorf("could not get new token")
	ErrUserNotFound      = fmt.Errorf("user not found")
)
