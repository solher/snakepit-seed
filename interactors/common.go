package interactors

import "github.com/solher/arangolite"

type QueryRunner interface {
	Run(q arangolite.Runnable) (*arangolite.Result, error)
}

type HTTPSender interface {
	Send(authPayload, method, url string, body interface{}) ([]byte, error)
}

type ConstantsGetter interface {
	GetString(key string) string
}
