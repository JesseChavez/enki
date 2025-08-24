package bouncer

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/JesseChavez/enki/cypher"
)

type Transcoder struct {
	secret        string
	salt          string
	maxAge        int
	derivedSecret []byte
}

func (tc *Transcoder) Init(secret string, salt string, maxAge int) *Transcoder {
	tc.secret = secret
	tc.salt   = salt
	tc.maxAge = maxAge

	key, err := cypher.KeyGenerator(secret, salt)

	if err != nil {
		panic(err)
	}

	tc.derivedSecret = key

	return  tc
}

func (tc *Transcoder) Decode(encodedValue string, decodedValue any) error {
	var err error

	// Format: output = [Data, IV, AuthTag]
	unpackedValue, err := tc.Unpack(encodedValue)

	if err != nil {
		return err
	}

	encMsg := unpackedValue[0]
	iv     := unpackedValue[1]
	tag    := unpackedValue[2]


	msg, err := cypher.DecryptMessage(tc.derivedSecret, encMsg, iv, tag)

	if err != nil {
		return err
	}

	value, err := UnwrapSession(msg)

	if err != nil {
		return err
	}

	err = json.Unmarshal(value, decodedValue)

	return err
}

func (tc *Transcoder) Encode(decodedValue any) (string, error) {
	var err error

	encodedValue, err := json.Marshal(decodedValue)

	if err != nil {
		return "", err
	}

	rawMsg, err := WrapSession(encodedValue, tc.maxAge)

	if err != nil {
		return "", err
	}

	msg, iv, tag, err := cypher.EncryptMessage(tc.derivedSecret, rawMsg)

	if err != nil {
		return "", err
	}

	// Format: [Data, IV, AuthTag]
	parts := [][]byte{msg, iv, tag}

	packedValue := tc.Pack(parts)

	return packedValue, nil
}

func (tc *Transcoder) Unpack(packedData string) ([][]byte, error)  {
	var err error

	// Format: Data--IV--AuthTag
	rawParts := strings.Split(packedData, "--")

	parts := [][]byte{}

	if len(rawParts) != 3 {
		return parts, errors.New("Invalid length")
	}

	for _, item := range rawParts {
		part, err := base64.StdEncoding.DecodeString(item)
		if err != nil {
			log.Fatal(err)
		}

		parts = append(parts, part)
	}
	
	return parts, err
}

func (tc *Transcoder) Pack(rawParts [][]byte) string {
	parts := []string{}

	// Format: [Data, IV, AuthTag]
	for _, item := range rawParts {
		part := base64.StdEncoding.EncodeToString(item)

		parts = append(parts, part)
	}

	return strings.Join(parts, "--")
}

