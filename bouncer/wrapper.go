package bouncer

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"
)

type Payload struct {
	Message string `json:"message"`
	Expiry  string `json:"exp"`
	Purpose string `json:"pur"`
}

type CookieSession struct {
	 Payload  Payload `json:"_rails"`
}


func WrapSession (encodedValue []byte, maxAge int) ([]byte, error)  {
	message := base64.StdEncoding.EncodeToString(encodedValue)

	expiry := ExpireTime(maxAge)

	session := CookieSession{
		Payload: Payload{
			Message: message,
			Expiry: expiry,
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

func ExpireTime(minutes int) string {
	// time is formatted in iso8601
	now := time.Now().Add(time.Minute * time.Duration(minutes))

	return now.UTC().Format("2006-01-02T15:04:05Z07:00")
}
