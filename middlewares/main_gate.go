package middlewares

import (
	"net/http"

	"git.wid.la/versatile/versatile-server/errs"

	"github.com/ansel1/merry"
	"github.com/pressly/chi"
	"github.com/solher/snakepit"
	"golang.org/x/net/context"
)

type MainGate struct {
	json *snakepit.JSON
}

func NewMainGate(j *snakepit.JSON) func(next chi.Handler) chi.Handler {
	gate := &MainGate{json: j}
	return gate.middleware
}

func (c *MainGate) middleware(next chi.Handler) chi.Handler {
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

		next.ServeHTTPC(ctx, w, r)
	})
}
