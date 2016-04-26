package utils

import (
	"crypto/rand"

	"github.com/ansel1/merry"
	"github.com/solher/arangolite/filters"
	"github.com/solher/snakepit-seed/errs"
)

func FilterToAQL(tmpVar string, f *filters.Filter) (string, error) {
	filter, err := filters.ToAQL(tmpVar, f)
	if err != nil {
		return "", merry.Here(errs.InvalidFilter)
	}

	return filter, nil
}

func GenToken(strSize int) string {
	dictionary := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	var bytes = make([]byte, strSize)
	rand.Read(bytes)

	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}

	return string(bytes)
}
