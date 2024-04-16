package main

import (
	"github.com/ciazhar/go-zhar/pkg/paseto"
	"github.com/rs/zerolog/log"
)

func main() {

	data := map[string]interface{}{
		"key": "value",
	}

	keyString := "cb40eac36dd250498cbe842a026263a72bcc3a77f2c8e90aa476d72e69eafd30"
	token, err := paseto.CreateToken(data, keyString)
	if err != nil {
		return
	}

	claims, err := paseto.ParseToken(token, keyString)
	if err != nil {
		return
	}
	for s := range claims {
		log.Info().Msgf("%s: %s", s, claims[s])
	}
}
