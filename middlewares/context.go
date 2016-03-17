package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"sync"

	"golang.org/x/net/context"

	"git.wid.la/versatile/versatile-server/models"

	"github.com/palantir/stacktrace"
	"github.com/pressly/chi"
)

type contextKey int

const (
	contextUser contextKey = iota
	contextAccessToken
	contextSession
)

func UserFromCtx(ctx context.Context) (*models.User, error) {
	v := ctx.Value(contextUser)
	if v == nil {
		return nil, stacktrace.NewError("nil user in context")
	}
	return v.(*models.User), nil
}

func AccessTokenFromCtx(ctx context.Context) (string, error) {
	v := ctx.Value(contextAccessToken)
	if v == nil {
		return "", stacktrace.NewError("nil access token in context")
	}
	return v.(string), nil
}

func SessionFromCtx(ctx context.Context) (*models.Session, error) {
	v := ctx.Value(contextSession)
	if v == nil {
		return nil, stacktrace.NewError("nil session in context")
	}
	return v.(*models.Session), nil
}

func Context(next chi.Handler) chi.Handler {
	fn := func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			user    *models.User
			token   string
			session *models.Session
		)

		wg := sync.WaitGroup{}
		wg.Add(3)

		go func() {
			user = getCurrentUser(r)
			wg.Done()
		}()

		go func() {
			token = getAccessToken(r)
			wg.Done()
		}()

		go func() {
			session = getCurrentSession(r)
			wg.Done()
		}()

		wg.Wait()

		ctx = context.WithValue(ctx, contextUser, user)
		ctx = context.WithValue(ctx, contextAccessToken, token)
		ctx = context.WithValue(ctx, contextSession, session)

		next.ServeHTTPC(ctx, w, r)
	}

	return chi.HandlerFunc(fn)
}

func getCurrentUser(r *http.Request) *models.User {
	enc := r.Header.Get("Auth-Server-Payload")
	if enc == "" {
		return nil
	}

	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return nil
	}

	user := &models.User{}

	if err := json.Unmarshal(data, user); err != nil {
		return nil
	}

	return user
}

func getAccessToken(r *http.Request) string {
	token := ""

	if t := r.Header.Get("Auth-Server-Token"); t != "" {
		token = t
	}

	if t := r.URL.Query().Get("accessToken"); t != "" {
		token = t
	}

	return token
}

func getCurrentSession(r *http.Request) *models.Session {
	enc := r.Header.Get("Auth-Server-Session")
	if enc == "" {
		return nil
	}

	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		return nil
	}

	session := &models.Session{}

	if err := json.Unmarshal(data, session); err != nil {
		return nil
	}

	return session
}
