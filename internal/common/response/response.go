package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/example/gin-monolithic-boilerplate/internal/common/apperror"
	"github.com/example/gin-monolithic-boilerplate/internal/common/requestid"
)

type Envelope struct {
	OK        bool       `json:"ok"`
	Data      any        `json:"data,omitempty"`
	Error     *ErrorBody `json:"error,omitempty"`
	RequestID string     `json:"requestId,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func Success(c *gin.Context, status int, data any) {
	if status == 0 {
		status = http.StatusOK
	}

	c.JSON(status, Envelope{OK: true, Data: data, RequestID: requestid.Get(c)})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func Error(c *gin.Context, err error) {
	appErr := apperror.From(err)
	if appErr == nil {
		appErr = apperror.Internal()
	}

	c.JSON(appErr.Status, Envelope{
		OK: false,
		Error: &ErrorBody{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		},
		RequestID: requestid.Get(c),
	})
	c.Abort()
}
