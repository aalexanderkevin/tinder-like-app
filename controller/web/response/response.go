package response

import (
	"errors"
	"net/http"
	"strconv"
	"tinder-like-app/helper"

	"github.com/gin-gonic/gin"
)

type (
	ErrorResponse struct {
		Code     string      `json:"code,omitempty"`
		Error    string      `json:"error,omitempty"`
		Message  string      `json:"error_message,omitempty"`
		Payload  interface{} `json:"payload,omitempty"`
		HttpCode int         `json:"-"`
	}

	SuccessResponse struct {
		Success bool `json:"success" default:"true"`
	}

	Message struct {
		Success bool        `json:"success"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
)

const (
	ErrorDuplicateKey = "duplicate_key"
	ErrorNotFound     = "not_found"
	ErrorUnauthorized = "unauthorized"
	ErrorBadRequest   = "bad_request"
	ErrorServerError  = "server_error"
	ErrorForbidden    = "forbidden"
)

var (
	SuccessOK = SuccessResponse{
		Success: true,
	}

	ErrNotFound = ErrorResponse{
		Error:    ErrorNotFound,
		Message:  "Entry not found",
		HttpCode: http.StatusNotFound,
	}
	ErrBadRequest = ErrorResponse{
		Error:    ErrorBadRequest,
		Message:  "Bad request",
		HttpCode: http.StatusBadRequest,
	}
	ErrUnauthorized = ErrorResponse{
		Error:    ErrorUnauthorized,
		Message:  "Unauthorized, please login",
		HttpCode: http.StatusUnauthorized,
	}
	ErrForbidden = ErrorResponse{
		Error:    ErrorForbidden,
		Message:  "You are unauthorized for this request",
		HttpCode: http.StatusForbidden,
	}
	ErrDuplicate = ErrorResponse{
		Error:    ErrorDuplicateKey,
		Message:  "Created value already exists",
		HttpCode: http.StatusConflict,
	}
	ErrValidation = ErrorResponse{
		Error:    ErrorBadRequest,
		Message:  "Invalid parameters or payload",
		HttpCode: http.StatusUnprocessableEntity,
	}
	ErrServerError = ErrorResponse{
		Error:    ErrorServerError,
		Message:  "Something bad happened",
		HttpCode: http.StatusInternalServerError,
	}
)

func SendErrorResponse(c *gin.Context, err ErrorResponse, msg string) {
	ErrorWithPayload(c, err, msg, nil)
}

func ErrorWithPayload(c *gin.Context, err ErrorResponse, msg string, payload interface{}) {
	c.Writer.Header().Del("content-type")
	if msg != "" {
		err.Message = msg
	}
	if payload != nil {
		err.Payload = payload
	}
	status := http.StatusBadRequest
	if err.HttpCode != 0 {
		status = err.HttpCode
	}
	c.JSON(status, err)
}
func Success(c *gin.Context) {
	SuccessWithPayload(c, nil)
}

func Error(c *gin.Context, err error, payload interface{}) {
	var errType helper.Error
	if errors.As(err, &errType) {
		handleEntityError(c, errType, payload)
	} else {
		SendErrorResponse(c, ErrServerError, err.Error())
	}
}

func handleEntityError(c *gin.Context, err helper.Error, payload interface{}) {
	if err.Code == helper.ErrorDuplicate {
		c.JSON(http.StatusConflict, payload)
		return
	}

	SendErrorResponse(c, ErrorResponse{
		Code:     strconv.Itoa(err.Code),
		Message:  err.Error(),
		Error:    "",
		HttpCode: err.Code,
	}, err.Error())
}

func SuccessWithPayload(c *gin.Context, payload interface{}) {
	if payload != nil {
		c.JSON(http.StatusOK, payload)
		return
	}
	c.JSON(http.StatusOK, SuccessOK)
}
