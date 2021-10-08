package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strings"
)

type Token struct {
	AuthTime int64                  `json:"auth_time"`
	Issuer   string                 `json:"iss"`
	Audience string                 `json:"aud"`
	Expires  int64                  `json:"exp"`
	IssuedAt int64                  `json:"iat"`
	Subject  string                 `json:"sub,omitempty"`
	UID      string                 `json:"uid,omitempty"`
	Firebase FirebaseInfo           `json:"firebase"`
	Claims   map[string]interface{} `json:"-"`
}

type FirebaseInfo struct {
	SignInProvider string                 `json:"sign_in_provider"`
	Tenant         string                 `json:"tenant"`
	Identities     map[string]interface{} `json:"identities"`
}

type jwtHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
	KeyID     string `json:"kid,omitempty"`
}

func DecodeFirebaseIdToken(idToken string) (*Token, error) {
	var (
		header  jwtHeader
		payload Token
	)
	segments := strings.Split(idToken, ".")
	if err := decode(segments[0], &header); err != nil {
		return nil, err
	}

	if err := decode(segments[1], &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func decode(segment string, i interface{}) error {
	decoded, err := base64.RawURLEncoding.DecodeString(segment)
	if err != nil {
		return err
	}
	return json.NewDecoder(bytes.NewBuffer(decoded)).Decode(i)
}
