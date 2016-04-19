package app

import (
	"net/http"

	"gopkg.in/h2non/gentleman.v0"

	"github.com/solher/snakepit-seed/constants"
	"github.com/solher/snakepit-seed/handlers"
	"github.com/solher/snakepit-seed/middlewares"

	"github.com/Sirupsen/logrus"
	"github.com/pressly/chi"
	"github.com/solher/arangolite"
	"github.com/solher/snakepit"
	"github.com/spf13/viper"
)

func Builder(v *viper.Viper, l *logrus.Logger) http.Handler {
	router := chi.NewRouter()
	json := snakepit.NewJSON()
	db := arangolite.New().LoggerOptions(false, false, false)
	db.Connect(
		v.GetString(constants.DBURL),
		v.GetString(constants.DBName),
		v.GetString(constants.DBUserName),
		v.GetString(constants.DBUserPassword),
	)
	cli := gentleman.New()

	timer := snakepit.NewTimer("Middleware stack")

	router.Use(snakepit.NewSwagger())
	router.Use(snakepit.NewRequestID())
	router.Use(snakepit.NewLogger(l))
	router.Use(timer.Start)
	router.Use(snakepit.NewRecoverer(json))
	router.Use(middlewares.NewContext())
	router.Use(timer.End)

	router.Mount("/users", handlers.NewUsers(v, json, db, cli))

	return router
}
