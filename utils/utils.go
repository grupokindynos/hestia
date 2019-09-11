package utils

import (
	"encoding/json"
	"gopkg.in/square/go-jose.v2"
)

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func EncryptJWE(key string, payload interface{}) (string, error) {
	encrypter, _ := jose.NewEncrypter(jose.A128GCM, jose.Recipient{Algorithm: jose.PBES2_HS256_A128KW, Key: key}, nil)
	payloadBytes, _ := json.Marshal(payload)
	object, _ := encrypter.Encrypt(payloadBytes)
	return object.CompactSerialize()
}

func DecryptJWE(key string, token string) ([]byte, error) {
	object, err := jose.ParseEncrypted(token)
	if err != nil {
		return nil, err
	}
	return object.Decrypt(key)
}
