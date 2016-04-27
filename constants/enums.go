package constants

import "github.com/solher/snakepit-seed/models"

const (
	RoleAdmin, RoleDeveloper, RoleUser models.Role = "ADMIN", "DEVELOPER", "USER"
)

var Roles = []models.Role{RoleAdmin, RoleDeveloper, RoleUser}
