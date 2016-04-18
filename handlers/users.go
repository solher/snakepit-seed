package handlers

import (
	"net/http"

	"gopkg.in/h2non/gentleman.v0"

	"git.wid.la/versatile/versatile-server/controllers"
	"git.wid.la/versatile/versatile-server/errs"
	"git.wid.la/versatile/versatile-server/interactors"
	"git.wid.la/versatile/versatile-server/middlewares"
	"git.wid.la/versatile/versatile-server/repositories"
	"git.wid.la/versatile/versatile-server/validators"
	"github.com/pressly/chi"
	"github.com/solher/arangolite/filters"
	"github.com/solher/snakepit"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

type (
	UsersCtrl interface {
		Create(ctx context.Context, w http.ResponseWriter, r *http.Request)
		Find(ctx context.Context, w http.ResponseWriter, r *http.Request)
		Update(ctx context.Context, w http.ResponseWriter, r *http.Request)
		Delete(ctx context.Context, w http.ResponseWriter, r *http.Request)

		FindByKey(ctx context.Context, w http.ResponseWriter, r *http.Request)
		UpdateByKey(ctx context.Context, w http.ResponseWriter, r *http.Request)
		DeleteByKey(ctx context.Context, w http.ResponseWriter, r *http.Request)

		Signin(ctx context.Context, w http.ResponseWriter, r *http.Request)
		FindSelf(ctx context.Context, w http.ResponseWriter, r *http.Request)
		CurrentSession(ctx context.Context, w http.ResponseWriter, r *http.Request)
		Signout(ctx context.Context, w http.ResponseWriter, r *http.Request)
		UpdateSelfPassword(ctx context.Context, w http.ResponseWriter, r *http.Request)
	}

	Users struct {
		snakepit.Handler
		DB     DatabaseRunner
		Client *gentleman.Client
	}
)

func NewUsers(
	c *viper.Viper,
	j *snakepit.JSON,
	db DatabaseRunner,
	cli *gentleman.Client,
) func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h := &Users{
		Handler: *snakepit.NewHandler(c, j),
		DB:      db,
		Client:  cli,
	}
	return h.builder
}

func (h *Users) routes(c UsersCtrl) chi.Router {
	r := chi.NewRouter()

	// CRUD operations
	r.Post("/", c.Create)
	r.Get("/", c.Find)
	r.Put("/", c.Update)
	r.Delete("/", c.Delete)

	// CRUD by key operations
	r.Route("/:key", func(r chi.Router) {
		r.Get("/", c.FindByKey)
		r.Put("/", c.UpdateByKey)
		r.Delete("/", c.DeleteByKey)
	})

	// Custom routes
	r.Post("/signin", c.Signin)
	r.Route("/me", func(r chi.Router) {
		r.Get("/", c.FindSelf)
		r.Get("/session", c.CurrentSession)
		r.Post("/signout", c.Signout)
		r.Post("/password", c.UpdateSelfPassword)
	})

	return r
}

func (h *Users) builder(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	accessToken, _ := middlewares.GetAccessToken(ctx)
	currentUser, _ := middlewares.GetCurrentUser(ctx)
	currentSession, _ := middlewares.GetCurrentSession(ctx)

	urlParams := chi.URLParams(ctx)

	filter, err := filters.FromRequest(r)
	if err != nil {
		h.JSON.RenderError(ctx, w, http.StatusBadRequest, errs.APIFilterDecoding, err)
		return
	}

	context := controllers.UsersContext{
		AccessToken:    accessToken,
		CurrentUser:    currentUser,
		CurrentSession: currentSession,
		URLParams:      urlParams,
		Filter:         filter,
	}

	logger, _ := snakepit.GetLogger(ctx)

	repo := repositories.NewRepository(
		h.Constants,
		logger,
		h.JSON,
		h.DB,
		h.Client,
	)

	sessionsInter := interactors.NewSessions(
		h.Constants,
		logger,
		repo,
	)
	inter := interactors.NewUsers(
		h.Constants,
		logger,
		repo,
		sessionsInter,
	)

	validator := validators.NewUsers()

	ctrl := controllers.NewUsers(
		h.Constants,
		logger,
		h.JSON,
		context,
		inter,
		sessionsInter,
		validator,
	)

	h.routes(ctrl).ServeHTTPC(ctx, w, r)
}