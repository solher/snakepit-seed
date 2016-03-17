package cmd

import (
	"fmt"
	"os"

	"github.com/solher/snakepit/root"
	"github.com/solher/snakepit/run"
	"git.wid.la/versatile/versatile-server/constants"
)

func init() {
	root.Cmd.Use = "versatile-server"
	root.Cmd.Short = "A custom description !"

	root.Cmd.AddCommand(run.Cmd)

	// DATABASE
	root.Cmd.PersistentFlags().String("dbUrl", "", "database URL")
	root.Viper.BindPFlag(constants.DBURL, root.Cmd.PersistentFlags().Lookup("dbUrl"))
	root.Viper.RegisterAlias(constants.DBURL, "ARANGODB_PORT")

	root.Cmd.PersistentFlags().String("dbName", "conet", "database name")
	root.Viper.BindPFlag(constants.DBName, root.Cmd.PersistentFlags().Lookup("dbName"))

	root.Cmd.PersistentFlags().String("dbRootName", "root", "database root user name")
	root.Viper.BindPFlag(constants.DBRootName, root.Cmd.PersistentFlags().Lookup("dbRootName"))

	root.Cmd.PersistentFlags().String("dbRootPassword", "qwertyuiop", "database root user password")
	root.Viper.BindPFlag(constants.DBRootPassword, root.Cmd.PersistentFlags().Lookup("dbRootPassword"))

	root.Cmd.PersistentFlags().String("dbUserName", "conet", "database main user name")
	root.Viper.BindPFlag(constants.DBUserName, root.Cmd.PersistentFlags().Lookup("dbUserName"))

	root.Cmd.PersistentFlags().String("dbUserPassword", "qwertyuiop", "database main user password")
	root.Viper.BindPFlag(constants.DBUserPassword, root.Cmd.PersistentFlags().Lookup("dbUserPassword"))
}

func Execute() {
	if err := root.Cmd.Execute(); err != nil {
		fmt.Printf(err.Error())
		os.Exit(-1)
	}
}
