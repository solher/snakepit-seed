package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/solher/snakepit-seed/app"
	"github.com/solher/snakepit-seed/constants"
	"github.com/solher/snakepit/root"
	"github.com/solher/snakepit/run"
)

func init() {
	run.Builder = app.Builder

	run.Logger.Formatter = &logrus.TextFormatter{
		ForceColors: true,
	}

	// APP
	run.Cmd.PersistentFlags().String("policyName", "snakepit", "policy created when sign in")
	root.Viper.BindPFlag(constants.PolicyName, run.Cmd.PersistentFlags().Lookup("policyName"))

	// SERVICES
	run.Cmd.PersistentFlags().String("authServerUrl", "", "auth server URL")
	root.Viper.BindPFlag(constants.AuthServerURL, run.Cmd.PersistentFlags().Lookup("authServerUrl"))
	root.Viper.RegisterAlias(constants.AuthServerURL, "AUTH_SERVER_PORT")

	// SWAGGER
	run.Cmd.PersistentFlags().String("swaggerBasePath", "/", "Swagger base path")
	root.Viper.BindPFlag(constants.SwaggerBasePath, run.Cmd.PersistentFlags().Lookup("swaggerBasePath"))
	run.Cmd.PersistentFlags().String("swaggerScheme", "http", "Swagger scheme (http or https)")
	root.Viper.BindPFlag(constants.SwaggerScheme, run.Cmd.PersistentFlags().Lookup("swaggerScheme"))
}
