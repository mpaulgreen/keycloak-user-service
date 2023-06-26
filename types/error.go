package types

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"

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

func Conflict(detail string) Error {
	return Error{
		Code:   http.StatusConflict,
		Status: http.StatusText(http.StatusConflict),
		Detail: detail,
	}
}

func ErrorFromResponse(response *http.Response, errMessage string) Error {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error reading response body: %s", err))
	}
	_ = response.Body.Close()
	return Error{
		Code:   response.StatusCode,
		Status: response.Status,
		Detail: fmt.Sprintf("%s: %s", errMessage, string(body)),
	}
}

func (e Error) String() string {
	return fmt.Sprintf("%s [%d] %s", e.Status, e.Code, e.Detail)
}
