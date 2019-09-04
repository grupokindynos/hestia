package controllers

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type FirebaseController struct {
	App *firebase.App
	DB  *mongo.Database
}

func (fb *FirebaseController) CheckAuth(c *gin.Context) {
	token := c.GetHeader("token")
	//payload := c.GetHeader("payload")

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
	_ = tk.UID
	/*// Decrypt the token payload
	encryptedToken, err := jose.ParseEncrypted(payload)
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorUnableParseJWE, c)
		return
	}
	decryptedToken, err := encryptedToken.Decrypt([]byte(uid))
	if err != nil {
		config.GlobalResponseError(nil, config.ErrorUnableDecryptJWE, c)
		return
	}*/
	// Run the function
	/*	res, err := function()
		if err != nil {
			config.GlobalResponseError(nil, err, c)
			return
		}*/
	config.GlobalResponseError(nil, nil, c)
	return
}

func (fb *FirebaseController) Return(uid string) (interface{}, error) {
	return "success " + uid, nil
}

func (fb *FirebaseController) ReturnAdmin() (interface{}, error) {
	return "success", nil
}
