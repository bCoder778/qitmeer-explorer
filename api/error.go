package api

import (
	"github.com/bCoder778/qitmeer-explorer/controller"
)

const (
	ERROR_UNKNOWN                     = 1
	Error_Addr_Invalid                = 2
	ERROR_PARAM                       = 0x00010101
)


type Error struct {
	Code    int
	Message string
}

func ParseError(err error)*Error{
	code := ERROR_UNKNOWN
	switch err {
	case controller.InvalidAddr:
		code =  Error_Addr_Invalid
	}
	return &Error{
		Code:    code,
		Message: err.Error(),
	}
}