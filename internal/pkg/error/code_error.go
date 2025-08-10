package myerror

type ErrorCode string

const (
	// 400
	Error_InvalidRequest     ErrorCode = "E4000"
	Error_InvalidCreds       ErrorCode = "E4001"
	Error_InvalidToken       ErrorCode = "E4010"
	Error_UserNotLogin       ErrorCode = "E4011"
	Error_Unauthorized       ErrorCode = "E4030"
	Error_APILocked          ErrorCode = "E4031"
	Error_RoleNotAllowed     ErrorCode = "E4032"
	Error_RouteNotFound      ErrorCode = "E4040"
	Error_RecordNotFound     ErrorCode = "E4041"
	Error_RecordAlreadyExist ErrorCode = "E4090"

	// 500
	SystemError            ErrorCode = "E5000"
	SystemError_Database   ErrorCode = "E5001"
	SystemError_UploadFile ErrorCode = "E5002"
	SystemError_SendEmail  ErrorCode = "E5003"
)

// default message error each error code
var Message = map[ErrorCode]string{
	// 400
	Error_InvalidRequest: "Invalid request format or parameters",
	Error_InvalidCreds:   "Invalid email or password",

	// 401
	Error_InvalidToken: "Invalid or expired token",
	Error_UserNotLogin: "You must be logged in to access this resource",

	// 403
	Error_Unauthorized:   "You do not have permission to perform this action",
	Error_APILocked:      "This API is temporarily locked. Please try again later",
	Error_RoleNotAllowed: "Your role is not allowed to access this resource",

	// 404
	Error_RouteNotFound:  "Route not found",
	Error_RecordNotFound: "Record not found",

	// 409
	Error_RecordAlreadyExist: "Record already exists",

	// 500
	SystemError:            "Something went wrong",
	SystemError_Database:   "Something went wrong in database operation failed",
	SystemError_UploadFile: "Something went wrong in upload file",
	SystemError_SendEmail:  "Something went wrong in send email",
}

var HttpStatuses = map[ErrorCode]int{
	// 400 Bad Request
	Error_InvalidRequest: 400,
	Error_InvalidCreds:   400,

	// 401 Unauthorized
	Error_InvalidToken: 401,
	Error_UserNotLogin: 401,

	// 403 Forbidden
	Error_Unauthorized:   403,
	Error_APILocked:      403,
	Error_RoleNotAllowed: 403,

	// 404 Not Found
	Error_RouteNotFound:  404,
	Error_RecordNotFound: 404,

	// 409 Conflict
	Error_RecordAlreadyExist: 409,

	// 500 Internal Server Error
	SystemError:            500,
	SystemError_Database:   500,
	SystemError_UploadFile: 500,
	SystemError_SendEmail:  500,
}
