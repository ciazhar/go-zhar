package response

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type BaseResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"` // optional, bisa diisi http.StatusBadRequest, dsb
}

type DataResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"` // optional, bisa diisi http.StatusBadRequest, dsb
	Data    any    `json:"data"`
}

type ValidationError struct {
	Field      string `json:"field,omitempty"` // empty for general errors
	Message    string `json:"message"`
	Suggestion string `json:"suggestion,omitempty"`
}

type ErrorResponse struct {
	Message string            `json:"message"`
	Code    int               `json:"code,omitempty"`   // optional, bisa diisi http.StatusBadRequest, dsb
	Errors  []ValidationError `json:"errors,omitempty"` // `omitempty` jika hanya satu message
}

type PageResponse struct {
	Code    int    `json:"code"` // optional, bisa diisi http.StatusBadRequest, dsb
	Message string `json:"message"`
	Data    any    `json:"data"`
	Total   int64  `json:"total"`
}

type CursorResponse struct {
	Code     int    `json:"code"` // optional, bisa diisi http.StatusBadRequest, dsb
	Message  string `json:"message"`
	Data     any    `json:"data"`
	PrevPath string `json:"prevPath,omitempty"`
	NextPath string `json:"nextPath,omitempty"`
	Total    int64  `json:"total"`
}

func NewBaseResponse(message string) BaseResponse {
	return BaseResponse{
		Message: message,
	}
}

func NewDataResponse(message string, data interface{}) DataResponse {
	return DataResponse{
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(message string, error error) ErrorResponse {
	return ErrorResponse{
		Message: message,
		Errors: []ValidationError{
			{
				Message: error.Error(),
			},
		},
	}
}

func NewErrorResponseWithErrors(message string, errors []ValidationError) ErrorResponse {
	return ErrorResponse{
		Message: message,
		Errors:  errors,
	}
}

func NewPageResponse(message string, data interface{}, total int64) PageResponse {
	return PageResponse{
		Message: message,
		Data:    data,
		Total:   total,
	}
}

// Success returns standardized success response
func Success(c *fiber.Ctx, code int, msg string, data interface{}) error {
	return c.Status(code).JSON(DataResponse{
		Message: msg,
		Data:    data,
	})
}

func PaginatedSuccess(c *fiber.Ctx, code int, msg string, data interface{}, total int64) error {
	return c.Status(code).JSON(PageResponse{
		Message: msg,
		Data:    data,
		Total:   total,
	})
}

func CursorSuccess(c *fiber.Ctx, code int, msg string, data interface{}, total int64, prevPath, nextPath string) error {
	return c.Status(code).JSON(CursorResponse{
		Message:  msg,
		Data:     data,
		PrevPath: prevPath,
		NextPath: nextPath,
		Total:    total,
	})
}

// HandleError logs and returns standardized error response
func HandleError(c *fiber.Ctx, log zerolog.Logger, code int, msg string, err error) error {
	log.Err(err).Msg(msg)
	return c.Status(code).JSON(ErrorResponse{
		Message: msg,
		Errors:  []ValidationError{{Message: err.Error()}},
	})
}

// HandleWarning logs and returns standardized warning response
func HandleWarning(c *fiber.Ctx, log zerolog.Logger, code int, msg string, err error) error {
	log.Warn().Err(err).Msg(msg)
	return c.Status(code).JSON(ErrorResponse{
		Message: msg,
		Errors:  []ValidationError{{Message: err.Error()}},
	})
}

// HandleWarnings logs and returns standardized warning response
func HandleWarnings(c *fiber.Ctx, log zerolog.Logger, code int, msg string, errs []ValidationError) error {
	log.Warn().Any("err", errs).Msg(msg)
	return c.Status(code).JSON(ErrorResponse{
		Message: msg,
		Errors:  errs,
	})
}

func ReturnError(log zerolog.Logger, msg string, err error) error {
	log.Err(err).Msg(msg)
	return err
}
