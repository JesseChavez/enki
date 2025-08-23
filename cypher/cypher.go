package cypher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

// NOTE: this package is called cypher rather cipher to avoid confusion with
// the crypo packages.

const (
	gcmTagSize = 16
)

func DecryptMessage(secret []byte, msg []byte, iv []byte, tag []byte) ([]byte, error) {
	// Go needs the auth tag appended to the data however Rails keeps
	// the auth tag separated from the data
	data := append(msg, tag...)

    // Setup cipher for decryption and add inputs
	aesCipher, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	mode, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, err
	}

	plaintext, err := mode.Open(nil, iv, data, nil)

	// fmt.Println("payload:", plaintext)
	return plaintext, err
}

func EncryptMessage(secret []byte, rawMsg []byte) ([]byte, []byte, []byte, error) {
    // Setup cipher for decryption and add inputs
	aesCipher, err := aes.NewCipher(secret)
	if err != nil {
		return nil, nil, nil, err
	}

	mode, err := cipher.NewGCM(aesCipher)
	if err != nil {
		return nil, nil, nil, err
	}

    iv := make([]byte, mode.NonceSize())

    _, err = rand.Read(iv)

    if err != nil {
		return nil, nil, nil, err
    }

	encMsg := mode.Seal(nil, iv, rawMsg, nil)

	msgLen := len(encMsg) - gcmTagSize

	// the auth tag is appended by GO to the data however Rails needs
	// the auth tag separated from the data
	msg, tag := encMsg[:msgLen], encMsg[msgLen:]

	return msg, iv, tag, nil
}
