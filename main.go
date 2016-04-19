//go:generate swagger generate spec -o swagger.json

// Snakepit boilerplate
//
// A simple Snakepit boilerplate.
//
// Schemes: https
// BasePath: /
// Version: 0.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package main // import "git.wid.la/versatile/versatile-server"
import "git.wid.la/versatile/versatile-server/cmd"

func main() {
	cmd.Execute()
}
