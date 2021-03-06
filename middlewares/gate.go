package middlewares

import (
	"net/http"

	"github.com/solher/snakepit-seed/constants"
	"github.com/solher/snakepit-seed/errs"
	"github.com/solher/snakepit-seed/models"

	"github.com/ansel1/merry"
	"github.com/pressly/chi"
	"github.com/solher/snakepit"
	"golang.org/x/net/context"
)

type Gate struct {
	json    *snakepit.JSON
	granter func(role models.Role) bool
}

func NewGate(j *snakepit.JSON, g func(role models.Role) bool) func(next chi.Handler) chi.Handler {
	gate := &Gate{json: j, granter: g}
	return gate.middleware
}

func (c *Gate) middleware(next chi.Handler) chi.Handler {
	return chi.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		session, err := GetCurrentSession(ctx)
		if err != nil {
			c.json.RenderError(ctx, w, http.StatusUnauthorized, errs.APIUnauthorized, err)
			return
		}
		if len(session.Role) == 0 {
			err := merry.New("empty role in context")
			c.json.RenderError(ctx, w, http.StatusUnauthorized, errs.APIUnauthorized, err)
			return
		}

		_, err = GetCurrentUser(ctx)
		if err != nil {
			c.json.RenderError(ctx, w, http.StatusUnauthorized, errs.APIUnauthorized, err)
			return
		}

		_, err = GetAccessToken(ctx)
		if err != nil {
			c.json.RenderError(ctx, w, http.StatusUnauthorized, errs.APIUnauthorized, err)
			return
		}

		if ok := c.granter(session.Role); !ok {
			err := merry.New("permission denied")
			c.json.RenderError(ctx, w, http.StatusForbidden, errs.APIForbidden, err)
			return
		}
		next.ServeHTTPC(ctx, w, r)
	})
}

func NewAdminOnly(j *snakepit.JSON) func(next chi.Handler) chi.Handler {
	gate := &Gate{
		json: j,
		granter: func(role models.Role) bool {
			return role == constants.RoleAdmin
		},
	}
	return gate.middleware
}

func NewAuthenticatedOnly(j *snakepit.JSON) func(next chi.Handler) chi.Handler {
	gate := &Gate{
		json: j,
		granter: func(role models.Role) bool {
			return len(role) != 0
		},
	}
	return gate.middleware
}
