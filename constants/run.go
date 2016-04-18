package constants

import "github.com/solher/snakepit/root"

const (
	AuthServerURL = "services.authServer.url"
)

const (
	RoleAdmin     = "roles.admin"
	RoleDeveloper = "roles.developer"
	RoleUser      = "roles.user"
)

func init() {
	root.Viper.Set(RoleAdmin, "ADMIN")
	root.Viper.Set(RoleDeveloper, "DEVELOPER")
	root.Viper.Set(RoleUser, "USER")
}
