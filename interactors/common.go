package interactors

import "github.com/solher/arangolite"

type QueryRunner interface {
	Run(q arangolite.Runnable, response interface{}) error
}

type HTTPSender interface {
	Send(authPayload, method, url string, body, response interface{}) error
}
