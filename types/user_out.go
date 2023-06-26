package types

import (
	"time"
)

type UserOut struct {
	Uuid string `json:"id"`

	Created time.Time `json:"created"`

	Modified time.Time `json:"modified"`

	UserId string `json:"user_id"`

	Username string `json:"username"`

	Email string `json:"email"`

	FirstName string `json:"first_name"`

	LastName string `json:"last_name"`

	OrgAdmin bool `json:"is_org_admin"`

	IsInternal bool `json:"is_internal"`

	OrgId string `json:"org_id"`

	Type_ string `json:"type"`

	IsActive bool `json:"is_active"` // This will be a manually calculated value

	// Added in addition to OpenAPI spec
	Enabled bool `json:"enabled"`
}
