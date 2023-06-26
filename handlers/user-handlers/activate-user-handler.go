package user_handles

import (
	"github.com/rs/zerolog/log"
	"keycloak-user-service/types"
	"strconv"
)

func (c *CallContext) ActivateUser(id string, activate bool) error {
	log.Info().Msg("Activating user with id: " + id + " with value: " + strconv.FormatBool(activate))

	// Get the user who needs to be activated
	user, err := c.client.GetUserByID(c.ctx, c.token, c.realm, id)
	if err != nil {
		log.Error().Msg("Error fetching user with id: " + id)
		return err
	}

	// Add approved custom attribute to the user if needed while activating
	if activate {
		attributes, err := c.effectiveAttributes(user)
		if err != nil {
			return err
		}
		approved, err := approvedAttribute(attributes)
		if err != nil {
			return err
		}
		if !approved {
			if user.Attributes == nil {
				blackAttrs := make(map[string][]string)
				user.Attributes = &blackAttrs
			}
			(*user.Attributes)[types.APPROVED_ATTRIBUTE_NAME] = []string{"true"}
		}
	}

	user.Enabled = &activate
	err = c.client.UpdateUser(c.ctx, c.token, types.KEYCLOAK_REALM, *user)
	if err != nil {
		log.Error().Msg("Cannot activate the user with id: " + id)
		return err
	}
	return nil
}
