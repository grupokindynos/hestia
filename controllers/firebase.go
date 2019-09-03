package controllers

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/config"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/square/go-jose.v2"
)

type FirebaseController struct {
	App *firebase.App
	DB  *mongo.Database
}

func (fb *FirebaseController) CheckAuth(c *gin.Context, function func(uid string, data []byte) (interface{}, error)) {
	token := c.GetHeader("token")
	payload := c.GetHeader("payload")

	// If there is no token on the header, return non-authed
	if token == "" {
		config.GlobalResponseNoAuth(c)
		return
	}
	// Verify token and get user information
	fbAuth, err := fb.App.Auth(context.Background())
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorFbInitializeAuth, c)
		return
	}
	tk, err := fbAuth.VerifyIDToken(context.Background(), token)
	if err != nil {
		config.GlobalResponseNoAuth(c)
		return
	}
	uid := tk.UID
	// Decrypt the token payload
	encryptedToken, err := jose.ParseEncrypted(payload)
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorUnableParseJWE, c)
		return
	}
	decryptedToken, err := encryptedToken.Decrypt([]byte(uid))
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorUnableDecryptJWE, c)
		return
	}
	// Run the function
	res, err := function(uid, decryptedToken)
	if err != nil {
		config.GlobalResponseError(nil, err, c)
		return
	}
	config.GlobalResponseError(res, nil, c)
	return
}

func (fb *FirebaseController) CheckAuthAdmin(c *gin.Context, method func(data []byte) (interface{}, error)) {
	token := c.GetHeader("token")
	payload := c.GetHeader("payload")

	// If there is no token on the header, return non-authed
	if token == "" {
		config.GlobalResponseNoAuth(c)
		return
	}
	// Verify token and get user information
	fbAuth, err := fb.App.Auth(context.Background())
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorFbInitializeAuth, c)
		return
	}
	tk, err := fbAuth.VerifyIDToken(context.Background(), token)
	if err != nil {
		config.GlobalResponseNoAuth(c)
		return
	}
	uid := tk.UID
	// TODO check for admin access

	// Decrypt the token payload
	encryptedToken, err := jose.ParseEncrypted(payload)
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorUnableParseJWE, c)
		return
	}
	decryptedToken, err := encryptedToken.Decrypt([]byte(uid))
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorUnableDecryptJWE, c)
		return
	}
	// Run the function
	res, err := method(decryptedToken)
	if err != nil {
		config.GlobalResponseError(nil, err, c)
		return
	}
	config.GlobalResponseError(res, nil, c)
	return
}

func (fb *FirebaseController) Return(uid string, data []byte) (interface{}, error) {
	return "success " + uid, nil
}

func (fb *FirebaseController) ReturnAdmin(data []byte) (interface{}, error) {
	return "success", nil
}
