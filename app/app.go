package app

import (
	"net/http"

	"git.wid.la/versatile/versatile-server/interactors"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/solher/arangolite"
	"github.com/solher/snakepit"
	"git.wid.la/versatile/versatile-server/controllers"
	"git.wid.la/versatile/versatile-server/repositories"
	"github.com/spf13/viper"
)

func Builder(v *viper.Viper) http.Handler {
	router := chi.NewRouter()
	render := snakepit.NewRender()
	db := arangolite.New()

	repository := repositories.NewRepository(db)

	sessionsInter := interactors.NewSessionsInter(v, repository)
	usersInter := interactors.NewUsersInter(repository, nil, sessionsInter)

	usersCtrl := controllers.NewUsersCtrl(usersInter, sessionsInter, nil, render)
	dashboardsCtrl := controllers.NewDashboardsCtrl(nil, nil, render)

	AddMiddlewares(router)
	AddRoutes(router, usersCtrl, dashboardsCtrl)

	return router
}

func AddMiddlewares(r chi.Router) {
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
}
