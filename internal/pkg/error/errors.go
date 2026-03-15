package myerror

var (
	ErrGeneral     = New("something went wrong", SystemError)
	ErrBodyRequest = New("failed get body request", Error_InvalidRequest)
)

type Error struct {
	Message    string
	ErrorCode  ErrorCode
	Detail     string
	StatusCode int
}

func New(msg string, statusCode ErrorCode) Error {
	return Error{
		Message:    msg,
		ErrorCode:  statusCode,
		StatusCode: HttpStatuses[statusCode],
	}
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) Details() string {
	return e.Detail
}
