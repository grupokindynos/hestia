package main

import (
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/hestia/config"
	"github.com/grupokindynos/hestia/controllers"
	"github.com/grupokindynos/hestia/models"
	"github.com/grupokindynos/hestia/services"
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
		// Init Database
		db, err := config.ConnectDB()
		if err != nil {
			log.Fatal(config.ErrorDbInitialize)
		}

		// Init Services
		obol := &services.ObolService{URL: "https://obol-rates.herokuapp.com/complex"}

		// Init DB models
		shiftsModel := &models.ShiftModel{Db: db}
		cardsModel := &models.CardsModel{Db: db}
		ordersModel := &models.OrdersModel{Db: db}
		depositsModel := &models.DepositsModel{Db: db}
		vouchersModel := &models.VouchersModel{Db: db}
		usersModel := &models.UsersModel{Db: db}

		// Init Controllers
		fbCtrl := controllers.FirebaseController{App: fbApp}
		_ = controllers.CardsController{Model: cardsModel}
		_ = controllers.DepositsController{Model: depositsModel}
		_ = controllers.OrdersController{Model: ordersModel}
		_ = controllers.ShiftsController{Model: shiftsModel, Obol: obol}
		_ = controllers.UsersController{Model: usersModel}
		_ = controllers.VouchersController{Model: vouchersModel}

		// Shift

		// User
		api.GET("/user/shift/single/:shiftid", func(c *gin.Context) { fbCtrl.CheckAuth(c) })
		api.GET("/user/shift/all", func(c *gin.Context) { fbCtrl.CheckAuth(c) })
		api.POST("/user/shift/new", func(c *gin.Context) { fbCtrl.CheckAuth(c) })

		// Admin
		api.GET("/shift/single/:shiftid", func(c *gin.Context) { fbCtrl.CheckAuth(c) })
		api.GET("/shift/all", func(c *gin.Context) { fbCtrl.CheckAuth(c) })

		// Vouchers

		// User
		api.GET("/user/voucher/single/:voucherid", func(c *gin.Context) { fbCtrl.CheckAuth(c) })
		api.GET("/user/voucher/all", func(c *gin.Context) { fbCtrl.CheckAuth(c) })
		api.POST("/user/voucher/new", func(c *gin.Context) { fbCtrl.CheckAuth(c) })

		// Admin
		api.GET("/voucher/single/:voucherid", func(c *gin.Context) { fbCtrl.CheckAuth(c) })
		api.GET("/voucher/all", func(c *gin.Context) { fbCtrl.CheckAuth(c) })

		// Cards

		// User
		api.GET("/user/cards/single/:cardcode", func(c *gin.Context) { fbCtrl.CheckAuth(c) })
		api.GET("/user/cards/all", func(c *gin.Context) { fbCtrl.CheckAuth(c) })

		// Admin
		api.GET("/cards/single/:cardcode", func(c *gin.Context) { fbCtrl.CheckAuth(c) })
		api.GET("/cards/all", func(c *gin.Context) { fbCtrl.CheckAuth(c) })
		api.POST("/cards/new", func(c *gin.Context) { fbCtrl.CheckAuth(c) })

		// Deposits

		// User
		api.GET("/user/deposits/single/:depositid", func(c *gin.Context) { fbCtrl.CheckAuth(c) })
		api.GET("/user/deposits/all", func(c *gin.Context) { fbCtrl.CheckAuth(c) })
		api.POST("/user/deposits/new", func(c *gin.Context) { fbCtrl.CheckAuth(c) })

		// Admin
		api.GET("/deposits/single/:depositid", func(c *gin.Context) { fbCtrl.CheckAuth(c) })
		api.GET("/deposits/all", func(c *gin.Context) { fbCtrl.CheckAuth(c) })

		// Orders

		// User
		api.GET("/user/orders/single/:orderid", func(c *gin.Context) { fbCtrl.CheckAuth(c) })
		api.GET("/user/orders/all", func(c *gin.Context) { fbCtrl.CheckAuth(c) })
		api.GET("/user/orders/new", func(c *gin.Context) { fbCtrl.CheckAuth(c) })

		// Admin
		api.GET("/orders/single/:orderid", func(c *gin.Context) { fbCtrl.CheckAuth(c) })
		api.GET("/orders/all", func(c *gin.Context) { fbCtrl.CheckAuth(c) })

		// User Information

		// User
		api.GET("/user/info", func(c *gin.Context) { fbCtrl.CheckAuth(c) })
		api.GET("/user/update", func(c *gin.Context) { fbCtrl.CheckAuth(c) })

		// Admin
		api.GET("/info/:uid", func(c *gin.Context) { fbCtrl.CheckAuth(c) })
		api.POST("/info/:uid", func(c *gin.Context) { fbCtrl.CheckAuth(c) })
	}
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})
}
