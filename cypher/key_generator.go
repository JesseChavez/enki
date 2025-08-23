package cypher

import (
	"crypto/pbkdf2"
	"crypto/sha256"
)


func KeyGenerator(password string, salt string) ([]byte, error) {
	// the password is the secret key base
	saltBytes := []byte(salt)

	// Choose a high number for security (current rails config)
	iterations := 1000

	// For a 256-bit key
	keyLength := 32

	// derivedKey, err := pbkdf2.Key(password, salt, iterations, keyLength, sha256.New)
	derivedKey, err := pbkdf2.Key(sha256.New, password, saltBytes, iterations, keyLength)

	// key := hex.EncodeToString(derivedKey)

	// fmt.Println("key:", key)

	return derivedKey, err
}
