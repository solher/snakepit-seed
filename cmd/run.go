package cmd

import (
	"github.com/solher/snakepit-seed/app"
	"github.com/solher/snakepit-seed/constants"
	"github.com/solher/snakepit/root"
	"github.com/solher/snakepit/run"
)

func init() {
	root.Cmd.AddCommand(run.Cmd)

	run.Builder = app.Builder

	// APP
	run.Cmd.PersistentFlags().String("policyName", "snakepit", "policy created when sign in")
	root.Viper.BindPFlag(constants.PolicyName, run.Cmd.PersistentFlags().Lookup("policyName"))

	// SERVICES
	run.Cmd.PersistentFlags().String("authServerUrl", "", "auth server URL")
	root.Viper.BindPFlag(constants.AuthServerURL, run.Cmd.PersistentFlags().Lookup("authServerUrl"))
	root.Viper.RegisterAlias(constants.AuthServerURL, "AUTH_SERVER_PORT")
}
