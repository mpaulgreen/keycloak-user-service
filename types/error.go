package types

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Error struct {
	Code   int    `json:"code,omitempty"`
	Status string `json:"status,omitempty"`
	Detail string `json:"detail,omitempty"`
}

func (err Error) Error() string {
	return err.String()
}

func InternalError(err error, errMessage string) Error {
	return Error{
		Code:   http.StatusInternalServerError,
		Status: http.StatusText(http.StatusInternalServerError),
		Detail: errors.Wrap(err, errMessage).Error(),
	}
}

func NotFound(detail string) Error {
	return Error{
		Code:   http.StatusNotFound,
		Status: http.StatusText(http.StatusNotFound),
		Detail: detail,
	}
}

func BadRequest(detail string) Error {
	return Error{
		Code:   http.StatusBadRequest,
		Status: http.StatusText(http.StatusBadRequest),
		Detail: detail,
	}
}

func (e Error) String() string {
	return fmt.Sprintf("%s [%d] %s", e.Status, e.Code, e.Detail)
}

type UserServiceErrors interface {
	CustomerGroupNotFoundError() error
	MoreThanOneCustomerGroup() error
	InternalServerError() error
}
