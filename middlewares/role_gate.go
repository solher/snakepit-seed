package middlewares

import (
	"net/http"

	"git.wid.la/versatile/versatile-server/errs"

	"github.com/ansel1/merry"
	"github.com/pressly/chi"
	"github.com/solher/snakepit"
	"golang.org/x/net/context"
)

type Role string

const (
	Admin, Developer, User Role = "ADMIN", "DEVELOPER", "USER"
)

type RoleGate struct {
	json    *snakepit.JSON
	granter func(role Role) bool
}

func NewRoleGate(j *snakepit.JSON, g func(role Role) bool) func(next chi.Handler) chi.Handler {
	gate := &RoleGate{json: j, granter: g}
	return gate.middleware
}

func (c *RoleGate) middleware(next chi.Handler) chi.Handler {
	return chi.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		session, _ := GetCurrentSession(ctx)

		if ok := c.granter(Role(session.Role)); !ok {
			err := merry.New("permission denied")
			c.json.RenderError(ctx, w, http.StatusForbidden, errs.APIForbidden, err)
			return
		}
		next.ServeHTTPC(ctx, w, r)
	})
}

func NewAdminOnly(j *snakepit.JSON) func(next chi.Handler) chi.Handler {
	gate := &RoleGate{
		json: j,
		granter: func(role Role) bool {
			return role == Admin
		},
	}
	return gate.middleware
}
