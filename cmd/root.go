package cmd

import (
	"github.com/solher/snakepit-seed/constants"
	"github.com/solher/snakepit/root"
)

func init() {
	root.Cmd.Use = "snakepit-seed"
	root.Cmd.Short = "A simple Snakepit seed."

	// DATABASE
	root.Cmd.PersistentFlags().String("dbUrl", "http://localhost:8000", "database URL")
	root.Viper.BindPFlag(constants.DBURL, root.Cmd.PersistentFlags().Lookup("dbUrl"))
	root.Viper.RegisterAlias(constants.DBURL, "ARANGODB_PORT")

	root.Cmd.PersistentFlags().String("dbName", "snakepit", "database name")
	root.Viper.BindPFlag(constants.DBName, root.Cmd.PersistentFlags().Lookup("dbName"))

	root.Cmd.PersistentFlags().String("dbRootName", "root", "database root user name")
	root.Viper.BindPFlag(constants.DBRootName, root.Cmd.PersistentFlags().Lookup("dbRootName"))

	root.Cmd.PersistentFlags().String("dbRootPassword", "qwertyuiop", "database root user password")
	root.Viper.BindPFlag(constants.DBRootPassword, root.Cmd.PersistentFlags().Lookup("dbRootPassword"))

	root.Cmd.PersistentFlags().String("dbUserName", "snakepit", "database main user name")
	root.Viper.BindPFlag(constants.DBUserName, root.Cmd.PersistentFlags().Lookup("dbUserName"))

	root.Cmd.PersistentFlags().String("dbUserPassword", "qwertyuiop", "database main user password")
	root.Viper.BindPFlag(constants.DBUserPassword, root.Cmd.PersistentFlags().Lookup("dbUserPassword"))
}

func Execute() {
	root.Execute()
}
