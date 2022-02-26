package api

import (
	"github.com/bCoder778/qitmeer-explorer/controller"
)

const (
	ERROR_UNKNOWN                     = 1
	Error_Addr_Invalid                = 2

	ERROR_REQUEST_UNAUTHPRIZED        = 0x00000191
	ERROR_REQUEST_NODFOUND            = 0x00000194
	ERROR_PARAM                       = 0x00010101
	ERROR_FORM_INIT_FAILED            = 0x00050101
)

type Error struct {
	Code    int 	`json:"code"`
	Message string  `json:"message"`
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