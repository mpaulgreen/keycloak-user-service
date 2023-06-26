package types

import (
	"fmt"
	"net/http"
)

type InviteUsers struct {
	Emails  []string `json:"emails"`
	IsAdmin bool     `json:"isAdmin"`
}

type InviteUsersResponse struct {
	Code    int
	Message string
}

func NewInviteUsersResponse() *InviteUsersResponse {
	return &InviteUsersResponse{Code: http.StatusCreated, Message: ""}
}

func (r *InviteUsersResponse) AddFailure(email string) {
	// TODO properly update the Code field. Use 207==StatusMultiStatus?
	r.Code = http.StatusPartialContent
	if len(r.Message) == 0 {
		r.Message = fmt.Sprintf("Completed with failures for users: %s", email)
	} else {
		r.Message = fmt.Sprintf("%s, %s", r.Message, email)
	}
}
