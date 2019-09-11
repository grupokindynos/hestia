package test

import (
	"encoding/json"
	"errors"
	"github.com/grupokindynos/hestia/utils"
	"github.com/stretchr/testify/assert"
	"gopkg.in/square/go-jose.v2"
	"testing"
)

func TestContains(t *testing.T) {
	array := []string{"a", "b"}
	contains := utils.Contains(array, "a")
	assert.Equal(t, true, contains)
}

func TestContainsError(t *testing.T) {
	array := []string{"a", "b"}
	contains := utils.Contains(array, "e")
	assert.Equal(t, false, contains)
}

func TestJWE(t *testing.T) {
	msg := "Encrypted Message"
	key := "RandomKey"
	encrypted, err := utils.EncryptJWE(key, msg)
	assert.Nil(t, err)
	decrypt, err := utils.DecryptJWE(key, encrypted)
	var msgStr string
	err = json.Unmarshal(decrypt, &msgStr)
	assert.Nil(t, err)
	assert.Equal(t, msg, msgStr)

}

func TestDecryptError(t *testing.T) {
	msg := "Encrypted Message"
	key := "RandomKey"
	encrypted, err := utils.EncryptJWE(key, msg)
	decrypt, err := utils.DecryptJWE("", encrypted)
	assert.Nil(t, decrypt)
	assert.NotNil(t, err)
	assert.Equal(t, jose.ErrCryptoFailure, err)
}

func TestDecryptInvalid(t *testing.T) {
	decrypt, err := utils.DecryptJWE("", "invalid-token")
	assert.Nil(t, decrypt)
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("square/go-jose: compact JWE format must have five parts"), err)
}
