package errors

import (
	"Technopark_DB_Project/app/models"
	"errors"
	"net/http"
)

var (
	// Forum errors
	ErrForumNotExist        = errors.New("Can't find user with id ") // TODO
	ErrForumOwnerNotFound   = errors.New("Can't find user with id ") // TODO
	ErrForumAlreadyExists   = errors.New("forum already exist")
	ErrForumOrTheadNotFound = errors.New("Can't find user with id ") // TODO

	// Thread errors
	ErrThreadAlreadyExists = errors.New("thread already exist")
	ErrThreadNotFound      = errors.New("Can't find user with id ") // TODO

	// Post errors
	ErrPostNotFound              = errors.New("Can't find user with id ") // TODO
	ErrParentPostNotExist        = errors.New("Can't find user with id ") // TODO
	ErrParentPostFromOtherThread = errors.New("Can't find user with id ")

	// User errors
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserNotFound     = errors.New("Can't find user with id ") // TODO
	ErrUserDataConflict = errors.New("Can't find user with id ") // TODO

	// Request Errors
	ErrBadInputData = errors.New("bad input data")
	ErrBadRequest   = errors.New("bad request")

	// Internal errors
	ErrNotImplemented = errors.New("not implemented")
	ErrInternal       = errors.New("internal error")
)

var errorToCodeMap = map[error]int{
	// Forum errors
	ErrForumNotExist:        http.StatusNotFound,
	ErrForumOwnerNotFound:   http.StatusNotFound,
	ErrForumAlreadyExists:   http.StatusConflict,
	ErrForumOrTheadNotFound: http.StatusNotFound,
	ErrThreadAlreadyExists:  http.StatusConflict,

	// Thread errors
	ErrThreadAlreadyExists: http.StatusConflict,
	ErrThreadNotFound:      http.StatusNotFound,

	// Post errors
	ErrPostNotFound:              http.StatusNotFound,
	ErrParentPostNotExist:        http.StatusNotFound,
	ErrParentPostFromOtherThread: http.StatusConflict,

	// User errors
	ErrUserAlreadyExist: http.StatusConflict,
	ErrUserNotFound:     http.StatusNotFound,
	ErrUserDataConflict: http.StatusConflict,

	// Request errors
	ErrBadInputData: http.StatusBadRequest,
	ErrBadRequest:   http.StatusBadRequest,

	// Internal errors
	ErrNotImplemented: http.StatusNotImplemented,
	ErrInternal:       http.StatusInternalServerError,
}

func ResolveErrorToCode(err error) (code int) {
	code, isErrorFound := errorToCodeMap[err]
	if !isErrorFound {
		code = http.StatusInternalServerError
	}
	return
}

func PrepareErrorResponse(err error) (statusCode int, contentType string, errorJSON []byte) {
	statusCode = ResolveErrorToCode(err)
	contentType = "application/json; charset=utf-8"
	errorJSON, errMarshal := models.Error{Message: err.Error()}.MarshalJSON()
	if errMarshal != nil {
		statusCode = ResolveErrorToCode(ErrInternal)
		errorJSON, _ = models.Error{Message: ErrInternal.Error()}.MarshalJSON()
	}
	return
}
