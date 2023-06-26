package user_handles

import (
	"fmt"
	"keycloak-user-service/types"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/Nerzal/gocloak/v13"

	"github.com/rs/zerolog/log"
)

func (c *CallContext) FindUsers(criteria types.FindUsersCriteria) (pagination types.UserPagination, code int, err error) {
	var usersList []types.UserOut
	if criteria.OrgId != nil {
		usersList, err = c.findByOrgId(criteria)
	} else if len(criteria.Emails()) > 0 || len(criteria.UserIds()) > 0 {
		err = errors.New("Queries by email or userId without an OrgId are not supported")
	} else if len(criteria.Usernames()) > 0 {
		usersList, err = c.findByUsernames(criteria.Usernames())
	} else {
		err = errors.New("Retrieval of all users is not supported")
	}
	if err != nil {
		log.Err(err).Msg("Failed to find users")
		code = http.StatusInternalServerError
		return
	}

	usersList = sortUsersList(criteria, usersList)
	paginationMeta := getPaginationObject(criteria, usersList)
	pagination = getPagedResults(criteria, usersList, paginationMeta)
	code = http.StatusOK
	return
}

func sortUsersList(findUsersCriteria types.FindUsersCriteria, usersList []types.UserOut) []types.UserOut {
	switch findUsersCriteria.OrderBy {
	case types.ORDER_BY_EMAIL:
		if findUsersCriteria.OrderDirection == types.ORDER_BY_DIR_ASC {
			return SortByEmail(usersList, true)
		} else {
			return SortByEmail(usersList, false)
		}
	case types.ORDER_BY_USERNAME:
		if findUsersCriteria.OrderDirection == types.ORDER_BY_DIR_ASC {
			return SortByUserName(usersList, true)
		} else {
			return SortByUserName(usersList, false)
		}
	case types.ORDER_BY_CREATED:
		if findUsersCriteria.OrderDirection == types.ORDER_BY_DIR_ASC {
			return SortByCreatedAt(usersList, true)
		} else {
			return SortByCreatedAt(usersList, false)
		}
	case types.ORDER_BY_MODIFIED:
		if findUsersCriteria.OrderDirection == types.ORDER_BY_DIR_ASC {
			return SortByModifiedAt(usersList, true)
		} else {
			return SortByModifiedAt(usersList, false)
		}
	default:
		log.Debug().Msg("Invalid order by parameter for find users, ignoring sorting the results.")
	}

	return usersList
}
func (c *CallContext) findByUsernames(usernames []string) ([]types.UserOut, error) {
	var usersList []types.UserOut
	params := gocloak.GetUsersParams{
		BriefRepresentation: &types.FALSE,
		Exact:               &types.TRUE,
	}
	for _, username := range usernames {
		//Optimizing by using keycloak queries for each given username
		params.Username = &username
		users, err := c.client.GetUsers(c.ctx, c.token, c.realm, params)
		if err != nil {
			return nil, err
		}
		for _, user := range users {
			attrs, err := c.effectiveAttributes(user)
			if err != nil {
				return nil, err
			}
			userOut, err := translate(user, attrs)
			if err != nil {
				return nil, err
			}
			usersList = append(usersList, *userOut)
		}
	}
	return usersList, nil
}

func (c *CallContext) findByUserID(userId string) ([]types.UserOut, error) {
	var usersList []types.UserOut
	params := gocloak.GetUsersParams{
		BriefRepresentation: &types.FALSE,
		Exact:               &types.TRUE,
	}

	//Optimizing by using keycloak queries for each given username
	params.Username = &userId
	users, err := c.client.GetUsers(c.ctx, c.token, c.realm, params)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		attrs, err := c.effectiveAttributes(user)
		if err != nil {
			return nil, err
		}
		userOut, err := translate(user, attrs)
		if err != nil {
			return nil, err
		}
		usersList = append(usersList, *userOut)
	}
	return usersList, nil
}

func (c *CallContext) populateHeritage(group types.GroupWrapper) error {
	//Unlike ancestry, group heritage is pre-populated so no server calls are needed to populate attribute info
	subgroups := group.Group().SubGroups
	if subgroups != nil && len(*subgroups) > 0 {
		for index := range *subgroups {
			wrappedSubgroup := group.AddChild(&(*subgroups)[index])
			if err := c.populateHeritage(wrappedSubgroup); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *CallContext) findByOrgId(criteria types.FindUsersCriteria) ([]types.UserOut, error) {
	//User query approach needs to be based on finding the group, getting all nested member users, and then filtering in code
	params := gocloak.GetGroupsParams{
		Q:                   gocloak.StringP(fmt.Sprintf("%s:%d", types.ORG_ID_ATTRIBUTE, *criteria.OrgId)),
		Full:                gocloak.BoolP(true),
		BriefRepresentation: gocloak.BoolP(false),
	}
	userMap := make(map[string]*types.UserOut)
	groups, err := c.client.GetGroups(c.ctx, c.token, c.realm, params)
	if err != nil {
		log.Err(err).Msg("Error getting groups for provided Org Id")
		return nil, err
	}
	log.Debug().Msg(fmt.Sprintf("Found %d groups with the provided organization ID", len(groups)))
	var matchedGroups []types.GroupWrapper
	for _, group := range groups {
		loadedGroup := types.WrapGroup(group)
		err = c.populateHeritage(loadedGroup)
		if err != nil {
			return nil, err
		}
		matchedGroups = append(matchedGroups, loadedGroup)
	}
	for _, group := range matchedGroups {
		err = c.populateGroupMembers(group, userMap)
		if err != nil {
			return nil, err
		}
	}
	var users []types.UserOut
	for _, user := range userMap {
		users = append(users, *user)
	}
	users = locallyFilterUsers(criteria, users)
	return users, nil
}

func locallyFilterUsers(criteria types.FindUsersCriteria, users []types.UserOut) []types.UserOut {
	//Users retrieved through group/members cannot be filtered with a server side query, so need to locally filter results
	var filters []userFilter
	if hasContent(criteria.UserIds()) {
		filters = append(filters, userFilter{
			matcher: matchesUserId,
			filters: criteria.UserIds(),
		})
	}
	if hasContent(criteria.Usernames()) {
		filters = append(filters, userFilter{
			matcher: matchesUsername,
			filters: criteria.Usernames(),
		})
	}
	if hasContent(criteria.Emails()) {
		filters = append(filters, userFilter{
			matcher: matchesEmail,
			filters: criteria.Emails(),
		})
	}

	if len(filters) == 0 {
		return users
	}
	var filtered []types.UserOut
	for index, user := range users {
		if matchesAnyFilter(user, filters) {
			filtered = append(filtered, users[index])
		}
	}
	return filtered
}

func (c *CallContext) populateGroupMembers(group types.GroupWrapper, userMap map[string]*types.UserOut) error {
	users, err := c.client.GetGroupMembers(c.ctx, c.token, c.realm, *group.Group().ID, gocloak.GetGroupsParams{})
	if err != nil {
		return err
	}
	for _, user := range users {
		_, exists := userMap[*user.ID]
		if !exists {
			userMap[*user.ID], err = translate(user, group.InheritedAttributes())
			if err != nil {
				return err
			}
		}
	}
	for _, subGroup := range group.Children() {
		err = c.populateGroupMembers(subGroup, userMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func translate(user *gocloak.User, inheritedAttributes *map[string][]string) (*types.UserOut, error) {
	userOut := types.UserOut{
		Uuid:       handleNilString(user.ID),
		Created:    time.UnixMilli(*user.CreatedTimestamp),
		Modified:   time.UnixMilli(*user.CreatedTimestamp), //TODO where is this supposed to come from? Do we need it?
		UserId:     handleNilString(user.ID),
		Username:   handleNilString(user.Username),
		Email:      handleNilString(user.Email),
		FirstName:  handleNilString(user.FirstName),
		LastName:   handleNilString(user.LastName),
		OrgAdmin:   false,
		IsInternal: false,
		OrgId:      "",
		Type_:      "",
		IsActive:   false,
		Enabled:    handleNilBool(user.Enabled),
	}

	if user.Attributes != nil {
		for key, value := range *user.Attributes {
			(*inheritedAttributes)[key] = value
		}
	}

	orgId, err := orgIdAttribute(inheritedAttributes)
	if err != nil {
		return nil, err
	}
	userOut.OrgId = handleNilString(orgId)

	orgAdmin, err := orgAdminAttribute(inheritedAttributes)
	if err != nil {
		return nil, err
	}
	userOut.OrgAdmin = orgAdmin

	approved, err := approvedAttribute(inheritedAttributes)
	if err != nil {
		return nil, err
	}
	userOut.IsActive = approved && userOut.Enabled

	return &userOut, nil
}

func getPaginationObject(findUsersCriteria types.FindUsersCriteria, usersList []types.UserOut) types.PaginationMeta {
	totalUsers := len(usersList)
	pageSize := findUsersCriteria.QueryLimit
	currentIdx := findUsersCriteria.Offset

	first := ""
	previous := ""
	next := ""
	last := ""

	if totalUsers > pageSize && pageSize > 0 {
		if currentIdx > 0 {
			first = fmt.Sprintf("%s%d", "/users?offset=0&limit=", findUsersCriteria.QueryLimit)
		}

		previousIdx := currentIdx - pageSize
		if previousIdx >= 0 {
			previous = fmt.Sprintf("%s%d%s%d", "/users?offset=", previousIdx, "&limit=", findUsersCriteria.QueryLimit)
		}

		nextIdx := currentIdx + pageSize
		if nextIdx < totalUsers && nextIdx >= pageSize {
			next = fmt.Sprintf("%s%d%s%d", "/users?offset=", nextIdx, "&limit=", findUsersCriteria.QueryLimit)
		}

		lastIdx := totalUsers - (totalUsers % pageSize)
		if lastIdx < totalUsers && currentIdx != lastIdx {
			last = fmt.Sprintf("%s%d%s%d", "/users?offset=", lastIdx, "&limit=", findUsersCriteria.QueryLimit)
		} else if lastIdx == totalUsers {
			last = fmt.Sprintf("%s%d%s%d", "/users?offset=", lastIdx-1, "&limit=", findUsersCriteria.QueryLimit)
		}
	}

	paginationMeta := types.PaginationMeta{
		Total:    int64(len(usersList)),
		First:    first,
		Previous: previous,
		Next:     next,
		Last:     last,
	}

	log.Debug().Msg(fmt.Sprintf("FindUsers Pagination: %+v\n", paginationMeta))

	return paginationMeta
}

func getPagedResults(findUsersCriteria types.FindUsersCriteria, usersList []types.UserOut, paginationMeta types.PaginationMeta) types.UserPagination {

	var returnUsersList []types.UserOut

	totalUsers := len(usersList)
	pageSize := findUsersCriteria.QueryLimit

	beginIdx := findUsersCriteria.Offset
	endIdx := beginIdx + pageSize
	if endIdx > totalUsers {
		endIdx = totalUsers
	}

	if beginIdx >= 0 && beginIdx < totalUsers && beginIdx <= endIdx {
		returnUsersList = usersList[beginIdx:endIdx]
	}

	userPagination := types.UserPagination{
		Meta:  &paginationMeta,
		Users: returnUsersList,
	}

	return userPagination
}
