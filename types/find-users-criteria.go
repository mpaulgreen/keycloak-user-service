package types

import "strings"

type FindUsersCriteria struct {
	OrgId               *int   `form:"org_id"`
	EmailsQueryArray    string `form:"emails"`
	UserIdsQueryArray   string `form:"user_ids"`
	UserNamesQueryArray string `form:"usernames"`
	QueryLimit          int    `form:"limit,default=1" binding:"omitempty,numeric,gte=1,max=1000"` // Max number of users to return
	Offset              int    `form:"offset,default=0" binding:"omitempty,numeric,gte=0"`
	OrderBy             string `form:"order" binding:"omitempty,oneof=email username modified created"` // values from specs: mail, username, modified, created
	OrderDirection      string `form:"direction" binding:"omitempty,oneof=asc desc"`                    // introduced by us, to be used only when Order parameter is specified. values: asc, desc
}

func (criteria *FindUsersCriteria) Emails() []string {
	var emails []string
	if len(criteria.EmailsQueryArray) != 0 {
		emails = strings.Split(criteria.EmailsQueryArray, ",")
	}
	return emails
}

func (criteria *FindUsersCriteria) Usernames() []string {
	var usernames []string
	if len(criteria.UserNamesQueryArray) != 0 {
		usernames = strings.Split(criteria.UserNamesQueryArray, ",")
	}
	return usernames
}

func (criteria *FindUsersCriteria) UserIds() []string {
	var userIds []string
	if len(criteria.UserIdsQueryArray) != 0 {
		userIds = strings.Split(criteria.UserIdsQueryArray, ",")
	}
	return userIds
}
