package bouncer

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

type Payload struct {
	Message string `json:"message"`
	Expiry  string `json:"exp"`
	Purpose string `json:"pur"`
}

type CookieSession struct {
	 Payload  Payload `json:"_rails"`
}


func WrapSession (encodedValue []byte) ([]byte, error)  {
	message := base64.StdEncoding.EncodeToString(encodedValue)

	session := CookieSession{
		Payload: Payload{
			Message: message,
		},
	}

	encodedSession, err := json.Marshal(session)

	return encodedSession, err
}

func UnwrapSession (msg []byte) ([]byte, error) {
	session := CookieSession{}

	err := json.Unmarshal(msg, &session)

	if err != nil {
		return nil, err
	}

	message := session.Payload.Message

	if message == "" {
		return nil, errors.New("Message attribute not found.")
	}

	decodedMessage, err := base64.StdEncoding.DecodeString(message)

	return decodedMessage, err
}
