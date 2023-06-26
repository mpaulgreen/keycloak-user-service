package user_handles

import (
	"keycloak-user-service/types"
	"strings"
)

type userFilter struct {
	matcher func(types.UserOut, string) bool
	filters []string
}

func matchesAnyFilter(user types.UserOut, userFilters []userFilter) bool {
	for _, filter := range userFilters {
		if matchesFilter(user, filter) {
			return true
		}
	}
	return false
}

func matchesFilter(user types.UserOut, usrFilter userFilter) bool {
	for _, filter := range usrFilter.filters {
		if usrFilter.matcher(user, filter) {
			return true
		}
	}
	return false
}

func matchesUserId(user types.UserOut, filter string) bool {
	return strings.Contains(user.UserId, filter)
}

func matchesUsername(user types.UserOut, filter string) bool {
	return strings.Contains(user.Username, filter)
}

func matchesEmail(user types.UserOut, filter string) bool {
	return strings.Contains(user.Email, filter)
}
