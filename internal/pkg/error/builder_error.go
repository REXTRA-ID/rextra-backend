package myerror

import "fmt"

func newError(code ErrorCode, message string, detail string) Error {
	if detail == "" {
		detail = message
	}

	return Error{
		Message:    message,
		ErrorCode:  code,
		StatusCode: HttpStatuses[code],
		Detail:     detail,
	}
}

func InvalidRequest(err error) Error {
	return newError(Error_InvalidRequest, Message[Error_InvalidRequest], err.Error())
}

func InvalidCreds() Error {
	return newError(Error_InvalidCreds, Message[Error_InvalidCreds], "")
}

func InvalidToken() Error {
	return newError(Error_InvalidToken, Message[Error_InvalidToken], "")
}

func Unauthorized() Error {
	return newError(Error_Unauthorized, Message[Error_Unauthorized], "")
}

func APILocked() Error {
	return newError(Error_APILocked, Message[Error_APILocked], "")
}

func RoleNotAllowed() Error {
	return newError(Error_RoleNotAllowed, Message[Error_RoleNotAllowed], "")
}

func UserNotLogin() Error {
	return newError(Error_UserNotLogin, Message[Error_UserNotLogin], "")
}

func RouteNotFound() Error {
	return newError(Error_RouteNotFound, Message[Error_RouteNotFound], "")
}

func RecordNotFound(item string) Error {
	msg := fmt.Sprintf("%s not found", item)
	return newError(Error_RecordNotFound, msg, msg)
}

func RecordAlreadyExist(item string) Error {
	msg := fmt.Sprintf("%s already exist", item)
	return newError(Error_RecordAlreadyExist, msg, msg)
}

func DatabaseError(err error) Error {
	return newError(SystemError_Database, Message[SystemError_Database], err.Error())
}

func UploadFileError(err error) Error {
	return newError(SystemError_UploadFile, Message[SystemError_UploadFile], err.Error())
}

func SendEmailError(err error) Error {
	return newError(SystemError_SendEmail, Message[SystemError_SendEmail], err.Error())
}
