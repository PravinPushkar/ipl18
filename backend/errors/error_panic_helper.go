package errors

import "log"

var PanicOnErr = func(err error, msg string) {
	if err != nil {
		log.Println(msg, err)
		panic(err)
	}
}
