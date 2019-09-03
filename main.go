package main

import (
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/controllers"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
)

func init() {
	_ = godotenv.Load()
}
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	App := GetApp()
	_ = App.Run(":" + port)
}

func GetApp() *gin.Engine {
	App := gin.Default()
	App.Use(cors.Default())
	fbCredStr := os.Getenv("FIREBASE_CRED")
	fbCred, err := base64.StdEncoding.DecodeString(fbCredStr)
	if err != nil {
		log.Fatal("unable to decode firebase credential string:")
	}
	opt := option.WithCredentialsJSON(fbCred)
	fbApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatal("unable to initialize firebase app")
	}
	ApplyRoutes(App, fbApp)
	return App
}

func ApplyRoutes(r *gin.Engine, fbApp *firebase.App) {
	api := r.Group("/")
	{
		fbCtrl := controllers.FirebaseController{App: fbApp}

		api.GET("/user/shift/single/:shiftid", func(c *gin.Context) { fbCtrl.CheckAuth(c, fbCtrl.Return) })
		api.GET("/user/shift/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, fbCtrl.Return) })
		api.GET("/user/voucher/single/:voucherid", func(c *gin.Context) { fbCtrl.CheckAuth(c, fbCtrl.Return) })
		api.GET("/user/voucher/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, fbCtrl.Return) })
		api.GET("/user/cards/single/:cardcode", func(c *gin.Context) { fbCtrl.CheckAuth(c, fbCtrl.Return) })
		api.GET("/user/cards/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, fbCtrl.Return) })
		api.GET("/user/deposits/single/:depositid", func(c *gin.Context) { fbCtrl.CheckAuth(c, fbCtrl.Return) })
		api.GET("/user/deposits/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, fbCtrl.Return) })
		api.GET("/user/orders/single/:orderid", func(c *gin.Context) { fbCtrl.CheckAuth(c, fbCtrl.Return) })
		api.GET("/user/orders/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, fbCtrl.Return) })
		api.GET("/user/info", func(c *gin.Context) { fbCtrl.CheckAuth(c, fbCtrl.Return) })
		api.GET("/shift/single/:shiftid", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, fbCtrl.ReturnAdmin) })
		api.GET("/shift/all", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, fbCtrl.ReturnAdmin) })
		api.GET("/voucher/single/:voucherid", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, fbCtrl.ReturnAdmin) })
		api.GET("/voucher/all", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, fbCtrl.ReturnAdmin) })
		api.GET("/cards/single/:cardcode", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, fbCtrl.ReturnAdmin) })
		api.GET("/cards/all", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, fbCtrl.ReturnAdmin) })
		api.GET("/deposits/single/:depositid", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, fbCtrl.ReturnAdmin) })
		api.GET("/deposits/all", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, fbCtrl.ReturnAdmin) })
		api.GET("/orders/single/:orderid", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, fbCtrl.ReturnAdmin) })
		api.GET("/orders/all", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, fbCtrl.ReturnAdmin) })
		api.GET("/info", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, fbCtrl.ReturnAdmin) })
	}
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})
}
