package utils

import (
	"crypto/rand"

	"git.wid.la/versatile/versatile-server/errs"
	"github.com/solher/arangolite/filters"
)

func FilterToAQL(tmpVar string, f *filters.Filter) (string, error) {
	filter, err := filters.ToAQL(tmpVar, f)
	if err != nil {
		return "", errs.InvalidFilter
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
