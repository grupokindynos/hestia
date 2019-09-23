package utils

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/jwt"
	"github.com/grupokindynos/hestia/config"
	"os"
)

func VerifyHeaderSignature(c *gin.Context) (string, error) {
	// Get the microservice name from header
	verificationToken := c.GetHeader("service")
	if verificationToken == "" {
		return "", errors.New("missing header signature")
	}
	// Decode not-verify
	payload, err := jwt.DecodeJWSNoVerify(verificationToken)
	// Unmarshal the payload
	var serviceStr string
	err = json.Unmarshal(payload, &serviceStr)
	if err != nil {
		return "", config.ErrorUnmarshal
	}
	// Check which service the request is announcing
	var pubKey string
	switch serviceStr {
	case "ladon":
		pubKey = os.Getenv("LADON_PUBLIC_KEY")
	case "tyche":
		pubKey = os.Getenv("TYCHE_PUBLIC_KEY")
	case "adrestia":
		pubKey = os.Getenv("ADRESTIA_PUBLIC_KEY")
	default:
		return "", config.ErrorNoAuth
	}
	// Verify the header token signature
	_, err = jwt.DecodeJWS(verificationToken, pubKey)
	// If there is an error, means the request was not properly signed
	if err != nil {
		return "", config.ErrorNoAuth
	}
	return pubKey, nil
}
