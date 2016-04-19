//go:generate swagger generate spec -o swagger.json

// Snakepit seed
//
// A simple Snakepit seed.
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
package main // import "github.com/solher/snakepit-seed"
import "github.com/solher/snakepit-seed/cmd"

func main() {
	cmd.Execute()
}
