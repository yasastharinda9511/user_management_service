package response

import "user_management_service/models"

type IntrospectResponse struct {
	Active bool         `json:"active"`
	User   *models.User `json:"user"`
}
