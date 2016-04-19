package cmd

import (
	"git.wid.la/versatile/versatile-server/app"
	"git.wid.la/versatile/versatile-server/constants"
	"github.com/solher/snakepit/root"
	"github.com/solher/snakepit/run"
)

func init() {
	run.Builder = app.Builder

	// APP
	run.Cmd.PersistentFlags().String("policyName", "snakepit", "policy created when sign in")
	root.Viper.BindPFlag(constants.PolicyName, run.Cmd.PersistentFlags().Lookup("policyName"))

	// SERVICES
	run.Cmd.PersistentFlags().String("authServerUrl", "", "auth server URL")
	root.Viper.BindPFlag(constants.AuthServerURL, run.Cmd.PersistentFlags().Lookup("authServerUrl"))
	root.Viper.RegisterAlias(constants.AuthServerURL, "AUTH_SERVER_PORT")
}
