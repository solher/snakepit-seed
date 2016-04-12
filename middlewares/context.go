package middlewares

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"golang.org/x/net/context"

	"git.wid.la/versatile/versatile-server/models"

	"github.com/Sirupsen/logrus"
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
		return nil, errors.New("nil context")
	}

	user, ok := ctx.Value(contextCurrentUser).(*models.User)
	if !ok {
		return nil, errors.New("unexpected type")
	}

	if user == nil {
		return nil, errors.New("nil value in context")
	}

	return user, nil
}

func GetAccessToken(ctx context.Context) (string, error) {
	if ctx == nil {
		return "", errors.New("nil context")
	}

	token, ok := ctx.Value(contextAccessToken).(string)
	if !ok {
		return "", errors.New("unexpected type")
	}

	if len(token) == 0 {
		return "", errors.New("empty value in context")
	}

	return token, nil
}

func GetCurrentSession(ctx context.Context) (*models.Session, error) {
	if ctx == nil {
		return nil, errors.New("nil context")
	}

	session, ok := ctx.Value(contextCurrentSession).(*models.Session)
	if !ok {
		return nil, errors.New("unexpected type")
	}

	if session == nil {
		return nil, errors.New("nil value in context")
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

		var (
			user    *models.User
			token   string
			session *models.Session
		)

		wg := sync.WaitGroup{}
		wg.Add(3)

		go func() {
			user = getCurrentUser(r, log)
			wg.Done()
		}()

		go func() {
			token = getAccessToken(r, log)
			wg.Done()
		}()

		go func() {
			session = getCurrentSession(r, log)
			wg.Done()
		}()

		wg.Wait()

		ctx = context.WithValue(ctx, contextCurrentUser, user)
		ctx = context.WithValue(ctx, contextAccessToken, token)
		ctx = context.WithValue(ctx, contextCurrentSession, session)

		next.ServeHTTPC(ctx, w, r)
	})
}

func getCurrentUser(r *http.Request, log *logrus.Entry) *models.User {
	enc := r.Header.Get("Auth-Server-Payload")
	if enc == "" {
		log.WithField("currentUser", enc).Debug("No current user set in context.")
		return nil
	}

	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		log.WithField("currentUser", enc).Debug("Could not base64 decode the current user in context.")
		return nil
	}

	user := &models.User{}

	if err := json.Unmarshal(data, user); err != nil {
		log.WithField("currentUser", data).Debug("Could not unmarshal the current user in context.")
		return nil
	}

	log.WithField("currentUser", data).Debug("Current user set in context.")

	return user
}

func getAccessToken(r *http.Request, log *logrus.Entry) string {
	token := r.Header.Get("Auth-Server-Token")

	if t := r.URL.Query().Get("accessToken"); t != "" {
		token = t
	}

	if len(token) == 0 {
		log.WithField("accessToken", "").Debug("No access token set in context.")
	} else {
		log.WithField("accessToken", token).Debug("Access token set in context.")
	}

	return token
}

func getCurrentSession(r *http.Request, log *logrus.Entry) *models.Session {
	enc := r.Header.Get("Auth-Server-Session")
	if enc == "" {
		log.WithField("currentSession", enc).Debug("No current session set in context.")
		return nil
	}

	data, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		log.WithField("currentSession", enc).Debug("Could not base64 decode the current session in context.")
		return nil
	}

	session := &models.Session{}

	if err := json.Unmarshal(data, session); err != nil {
		log.WithField("currentSession", data).Debug("Could not unmarshal the current session in context.")
		return nil
	}

	log.WithField("currentSession", data).Debug("Current session set in context.")

	return session
}
