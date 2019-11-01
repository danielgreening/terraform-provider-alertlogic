package al_client_base

import (
	"strings"

	"algithub.pd.alertlogic.net/daniel-greening/alertlogic-sdk-go/alertlogic/alerr"
)

// IsALErr returns true if the error matches all these conditions:
//  * err is of type alerr.Error
//  * Error.Code() matches code
//  * Error.Message() contains message
func IsALErr(err error, code string, message string) bool {
	alErr, ok := err.(alerr.Error)

	if !ok {
		return false
	}

	if alErr.Code() != code {
		return false
	}

	return strings.Contains(alErr.Message(), message)
}

// IsALsErrExtended returns true if the error matches all these conditions:
//  * err is of type alerr.Error
//  * Error.Code() matches code
//  * Error.Message() contains message
//  * Error.OrigErr() contains origErrMessage
func IsALErrExtended(err error, code string, message string, origErrMessage string) bool {
	if !IsALErr(err, code, message) {
		return false
	}
	return strings.Contains(err.(alerr.Error).OrigErr().Error(), origErrMessage)
}
