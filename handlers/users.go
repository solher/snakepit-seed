package handlers

import (
	"net/http"

	"git.wid.la/versatile/versatile-server/controllers"
	"git.wid.la/versatile/versatile-server/interactors"
	"git.wid.la/versatile/versatile-server/repositories"
	"github.com/pressly/chi"
	"github.com/solher/arangolite"
	"github.com/solher/snakepit"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

func Users(
	vip *viper.Viper,
	render *snakepit.Render,
	db *arangolite.DB,
) func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		repo := repositories.New(db)

		sessionsInter := interactors.NewSessions(vip, repo)
		usersInter := interactors.NewUsers(repo, nil, sessionsInter)

		ctrl := controllers.NewUsers(usersInter, sessionsInter, nil, render)

		router := chi.NewRouter()

		// CRUD operations
		router.Route("/", func(router chi.Router) {
			router.Post("/", ctrl.Create)
			router.Get("/", ctrl.Find)
			router.Put("/", ctrl.Update)
			router.Delete("/", ctrl.Delete)
		})

		// CRUD by key operations
		router.Route("/:key", func(router chi.Router) {
			router.Get("/", ctrl.FindByKey)
			router.Put("/", ctrl.UpdateByKey)
			router.Delete("/", ctrl.DeleteByKey)
		})

		// Custom routes
		router.Post("/signin", ctrl.Signin)
		router.Route("/me", func(router chi.Router) {
			router.Get("/", ctrl.FindSelf)
			router.Get("/session", ctrl.CurrentSession)
			router.Post("/signout", ctrl.Signout)
			router.Post("/password", ctrl.UpdateSelfPassword)
		})

		router.ServeHTTPC(ctx, w, r)
	}
}
