package paseto

import (
	"aidanwoods.dev/go-paseto"
	"github.com/rs/zerolog/log"
	"time"
)

func CreateToken(data map[string]interface{}, keyString string) (string, error) {
	token := paseto.NewToken()

	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(time.Now().Add(2 * time.Hour))

	for s := range data {
		err := token.Set(s, data[s])
		if err != nil {
			return "", err
		}
	}

	key, err := paseto.V4SymmetricKeyFromHex(keyString) // don't share this!!
	if err != nil {
		return "", err
	}

	encrypted := token.V4Encrypt(key, nil)

	return encrypted, nil
}

func ParseToken(tokenString string, keyString string) (map[string]interface{}, error) {
	key, err := paseto.V4SymmetricKeyFromHex(keyString)
	if err != nil {
		return nil, err
	}

	parser := paseto.NewParser()
	token, err := parser.ParseV4Local(key, tokenString, nil)
	if err != nil {
		log.Error().Msg(err.Error())
	}

	if token == nil {
		return nil, err
	}

	return token.Claims(), nil
}
