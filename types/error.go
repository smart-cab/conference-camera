package types

type ApiError struct {
	Code int
	Msg  string
}

var (
	USER_EXIST          = &ApiError{202, "user already exist"}
	USER_NOT_FOUND      = &ApiError{202, "user does not exist"}
	WRONG_PASSWORD      = &ApiError{203, "wrong login or password"}
	WRONG_PARAMETERS    = &ApiError{300, "wrong parameters"}
	FAILED_DATABASE     = &ApiError{400, "operation failed"}
	UNAUTHORIZED        = &ApiError{401, "unauthorized"}
	NOT_FOUND           = &ApiError{404, "not found"}
	INCORRECT_WEBSOCKET = &ApiError{400, "invalid websocket"}
	ALREADY_CONNECTED   = &ApiError{403, "client already connected"}
)

func (e *ApiError) Error() string {
	return e.Msg
}
