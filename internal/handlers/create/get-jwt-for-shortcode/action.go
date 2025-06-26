package getJwtForShortCode

import (
	"encoding/base64"
	"fmt"
	"github.com/driscollco-core/encryption/keypair"
	"github.com/driscollco-core/jwt"
	serviceComponents "github.com/driscollco-core/service/components"
	"github.com/driscollco-shortify/entities"
	"github.com/driscollco-shortify/svc-shortcode/internal/conf"
)

func Action(bundle serviceComponents.Bundle, shortCode entities.ShortCode, isCreator bool) (string, error) {
	span := bundle.Span("decoding keys")
	public, err := decodeKey(conf.Config.Service.Handlers.Create.Jwt.PublicKey)
	if err != nil {
		span.Error("error decoding public key")
		return "", fmt.Errorf("could not decode public key: %w", err)
	}

	private, err := decodeKey(conf.Config.Service.Handlers.Create.Jwt.PrivateKey)
	if err != nil {
		span.Error("error decoding private key")
		return "", fmt.Errorf("could not decode private key: %w", err)
	}
	span.Success()

	span = bundle.Span("creating jwt for shortCode")
	pair, err := keypair.LoadKeypair(public, private)
	if err != nil {
		span.Error("error loading key pair")
		return "", fmt.Errorf("could not load public/private keypair : %w", err)
	}

	token := jwt.Token{}
	token.Claim("id", shortCode.Id)
	token.Claim("url", shortCode.ShortCode.URL)
	if isCreator {
		token.Claim("is.creator", true)
	}
	span.Success()
	return token.Create(pair.Private()), nil
}

func decodeKey(encoded string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 key: %w", err)
	}
	return string(decodedBytes), nil
}
