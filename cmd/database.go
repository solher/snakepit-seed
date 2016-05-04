package cmd

import (
	"strings"

	"github.com/solher/snakepit"
	"github.com/solher/snakepit-seed/constants"
	"github.com/solher/snakepit-seed/database"
	dbCmd "github.com/solher/snakepit/database"
	"github.com/solher/snakepit/root"
	"github.com/spf13/viper"
)

func init() {
	root.Cmd.AddCommand(dbCmd.Cmd)

	dbCmd.Create = func(v *viper.Viper) error {
		return initDatabaseManager(v).Create(
			v.GetString(constants.DBRootName),
			v.GetString(constants.DBRootPassword),
		)
	}

	dbCmd.Migrate = func(v *viper.Viper) error {
		return initDatabaseManager(v).Migrate()
	}

	dbCmd.Seed = func(v *viper.Viper) error {
		return initDatabaseManager(v).SyncSeeds()
	}

	dbCmd.Drop = func(v *viper.Viper) error {
		return initDatabaseManager(v).Drop(
			v.GetString(constants.DBRootName),
			v.GetString(constants.DBRootPassword),
		)
	}
}

func initDatabaseManager(v *viper.Viper) *database.Manager {
	ara := snakepit.NewArangoDBManager(database.NewProdSeed(), database.NewEmptyProdSeed()).
		LoggerOptions(false, false, false).
		Connect(
		strings.Replace(v.GetString(constants.DBURL), "tcp", "http", -1),
		v.GetString(constants.DBName),
		v.GetString(constants.DBUserName),
		v.GetString(constants.DBUserPassword),
	)

	return database.NewManager(ara)
}
