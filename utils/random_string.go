package utils

import "github.com/jaevor/go-nanoid"

func GenerateRandomString(len int) string {
	generated, _ := nanoid.Standard(len)
	return generated()
}
