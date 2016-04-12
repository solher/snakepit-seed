package app

import (
	"net/http"

	"git.wid.la/versatile/versatile-server/handlers"
	"git.wid.la/versatile/versatile-server/middlewares"

	"github.com/Sirupsen/logrus"
	"github.com/pressly/chi"
	"github.com/solher/arangolite"
	"github.com/solher/snakepit"
	"github.com/spf13/viper"
)

func Builder(v *viper.Viper, l *logrus.Logger) http.Handler {
	router := chi.NewRouter()
	render := snakepit.NewRender()
	db := arangolite.New()

	router.Use(snakepit.NewRequestID())
	router.Use(snakepit.NewLogger(l))
	router.Use(snakepit.NewRecoverer(render))
	router.Use(middlewares.NewContext())
	router.Handle("/users", handlers.Users(v, render, db))

	return router
}
