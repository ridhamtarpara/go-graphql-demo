package errors

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"runtime"
)

var (
	ServerError         = GenerateError("Something went wrong! Please try again later")
	UserNotExist        = GenerateError("User not exists")
	UnauthorisedError        = GenerateError("You are not authorised to perform this action")
	TimeStampError      = GenerateError("time should be a unix timestamp")
	InternalServerError = GenerateError("internal server error")
)

func GenerateError(err string) error {
	return errors.New(err)
}
func IsForeignKeyError(err error) bool {
	pgErr := err.(*pq.Error);
	if pgErr.Code == "23503" {
		return true
	}
	return false
}

func DebugPrintf(err_ error, args ...interface{}) string {
	programCounter, file, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(programCounter)
	msg := fmt.Sprintf("[%s: %s %d] %s, %s", file, fn.Name(), line, err_, args)
	log.Println(msg)
	return msg
}
