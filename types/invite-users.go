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

type InviteUsersError struct {
}

var (
	CUSTOMER_GROUP_NOT_FOUND     = fmt.Errorf("CUSTOMER_GROUP_NOT_FOUND")
	MORE_THAN_ONE_CUSTOMER_GROUP = fmt.Errorf("MORE_THAN_ONE_CUSTOMER_GROUP")
	INTERNAL_SERVER_ERROR        = fmt.Errorf("INTERNAL_SERVER_ERROR") // TODO: Place these and other common errors in error.go file
)

func (*InviteUsersError) CustomerGroupNotFoundError() error {
	return CUSTOMER_GROUP_NOT_FOUND
}

func (*InviteUsersError) MoreThanOneCustomerGroup() error {
	return MORE_THAN_ONE_CUSTOMER_GROUP
}

func (*InviteUsersError) InternalServerError() error {
	return INTERNAL_SERVER_ERROR
}
