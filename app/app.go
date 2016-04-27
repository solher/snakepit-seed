package app

import (
	"net/http"

	"gopkg.in/h2non/gentleman.v1"

	"github.com/solher/snakepit-seed/constants"
	"github.com/solher/snakepit-seed/database"
	"github.com/solher/snakepit-seed/handlers"
	"github.com/solher/snakepit-seed/middlewares"

	"github.com/Sirupsen/logrus"
	"github.com/pressly/chi"
	"github.com/solher/snakepit"
	"github.com/spf13/viper"
)

func Builder(v *viper.Viper, l *logrus.Logger) http.Handler {
	distantSeed := database.NewEmptyProdSeed()

	db := snakepit.NewArangoDBManager(
		database.NewProdSeed(),
		distantSeed,
	).
		LoggerOptions(false, false, false).
		Connect(
		v.GetString(constants.DBURL),
		v.GetString(constants.DBName),
		v.GetString(constants.DBUserName),
		v.GetString(constants.DBUserPassword),
	)

	if err := db.LoadDistantSeed(); err != nil {
		panic(err)
	}

	distantSeed.PopulateConstants(v)

	router := chi.NewRouter()
	json := snakepit.NewJSON()
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
