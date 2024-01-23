package utils

import (
	"math/rand"
)

func GenerateRandomString(length int) string {
	// Define the characters that can be used in the random string
	charset := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_="
	charsetLength := len(charset)

	// Create the random string
	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[rand.Intn(charsetLength)]
	}

	return string(randomString)
}
