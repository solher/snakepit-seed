package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"git.wid.la/versatile/versatile-server/models"

	"github.com/Sirupsen/logrus"
	"github.com/ansel1/merry"
	"github.com/pressly/chi"
	"github.com/solher/snakepit"
)

const (
	contextCurrentUser    snakepit.CtxKey = "currentUser"
	contextAccessToken    snakepit.CtxKey = "accessToken"
	contextCurrentSession snakepit.CtxKey = "currentSession"
)

func GetCurrentUser(ctx context.Context) (*models.User, error) {
	if ctx == nil {
		return nil, merry.New("nil context")
	}

	user, ok := ctx.Value(contextCurrentUser).(*models.User)
	if !ok {
		return nil, merry.New("unexpected type")
	}

	if user == nil {
		return nil, merry.New("nil value in context")
	}

	return user, nil
}

func GetAccessToken(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", merry.New("nil context")
	}

	token, ok := ctx.Value(contextAccessToken).(string)
	if !ok {
		return "", merry.New("unexpected type")
	}

	if len(token) == 0 {
		return "", merry.New("empty value in context")
	}

	return token, nil
}

func GetCurrentSession(ctx context.Context) (*models.Session, error) {
	if ctx == nil {
		return nil, merry.New("nil context")
	}

	session, ok := ctx.Value(contextCurrentSession).(*models.Session)
	if !ok {
		return nil, merry.New("unexpected type")
	}

	if session == nil {
		return nil, merry.New("nil value in context")
	}

	return session, nil
}

type Context struct{}

func NewContext() func(next chi.Handler) chi.Handler {
	context := &Context{}
	return context.middleware
}

func (c *Context) middleware(next chi.Handler) chi.Handler {
	return chi.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		log, _ := snakepit.GetLogger(ctx)

		payload := getAuthServerPayload(r, log)
		token := getAccessToken(r, log)
		session := getCurrentSession(r, log)
		if session != nil {
			session.Role = payload.Role
		}

		ctx = context.WithValue(ctx, contextCurrentUser, payload.User)
		ctx = context.WithValue(ctx, contextAccessToken, token)
		ctx = context.WithValue(ctx, contextCurrentSession, session)

		next.ServeHTTPC(ctx, w, r)
	})
}

func getAuthServerPayload(r *http.Request, log *logrus.Entry) *models.AuthServerPayload {
	payload := &models.AuthServerPayload{}

	enc := r.Header.Get("Auth-Server-Payload")
	if enc == "" {
		log.WithField("authServerPayload", enc).
			Debug("No auth server payload received.")
		return payload
	}

	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		log.WithField("authServerPayload", enc).
			Debug("Could not base64 decode the received auth server payload.")
		return payload
	}

	if err := json.Unmarshal(data, payload); err != nil {
		log.WithField("authServerPayload", string(data)).
			Debug("Could not unmarshal the received auth server payload.")
		return payload
	}

	log.WithField("authServerPayload", string(data)).
		Debug("Auth server payload received.")

	return payload
}

func getAccessToken(r *http.Request, log *logrus.Entry) string {
	token := r.Header.Get("Auth-Server-Token")

	if t := r.URL.Query().Get("accessToken"); t != "" {
		token = t
	}

	if len(token) == 0 {
		log.WithField("accessToken", "").
			Debug("No access token received.")
	} else {
		log.WithField("accessToken", token).
			Debug("Access token received.")
	}

	return token
}

func getCurrentSession(r *http.Request, log *logrus.Entry) *models.Session {
	enc := r.Header.Get("Auth-Server-Session")
	if enc == "" {
		log.WithField("currentSession", enc).
			Debug("No current session received.")
		return nil
	}

	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		log.WithField("currentSession", enc).
			Debug("Could not base64 decode the received current session.")
		return nil
	}

	session := &models.Session{}

	if err := json.Unmarshal(data, session); err != nil {
		log.WithField("currentSession", string(data)).
			Debug("Could not unmarshal the received current session.")
		return nil
	}

	log.WithField("currentSession", string(data)).
		Debug("Current session received.")

	return session
}
