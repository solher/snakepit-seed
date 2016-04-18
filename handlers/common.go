package handlers

import "github.com/solher/arangolite"

type DatabaseRunner interface {
	Run(q arangolite.Runnable) ([]byte, error)
}
