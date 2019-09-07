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
		_ = &services.PlutusService{URL: "https://obol-rates.herokuapp.com/complex", AuthUsername: os.Getenv("PLUTUS_AUTH_USERNAME"), AuthPassword: os.Getenv("PLUTUS_AUTH_PASSWORD")}

		// Init DB models
		shiftsModel := &models.ShiftModel{Db: db}
		cardsModel := &models.CardsModel{Db: db}
		ordersModel := &models.OrdersModel{Db: db}
		depositsModel := &models.DepositsModel{Db: db}
		vouchersModel := &models.VouchersModel{Db: db}
		usersModel := &models.UsersModel{Db: db}
		coinsModel := &models.CoinsModel{Db: db}

		// Init Controllers
		fbCtrl := controllers.FirebaseController{App: fbApp}
		cardsCtrl := controllers.CardsController{Model: cardsModel}
		depositsCtrl := controllers.DepositsController{Model: depositsModel}
		ordersCtrl := controllers.OrdersController{Model: ordersModel}
		shiftCtrl := controllers.ShiftsController{Model: shiftsModel, Obol: obol}
		userCtrl := controllers.UsersController{Model: usersModel}
		vouchersCtrl := controllers.VouchersController{Model: vouchersModel}
		coinsCtrl := controllers.CoinsController{Model: coinsModel}

		// User Information

		// User
		api.GET("/user/info", func(c *gin.Context) { fbCtrl.CheckAuth(c, userCtrl.GetSelfInfo) })
		api.GET("/user/update", func(c *gin.Context) { fbCtrl.CheckAuth(c, userCtrl.UpdateSelfInfo) })

		// Admin
		api.GET("/info/:uid", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, userCtrl.GetUserInfo) })
		api.POST("/info", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, userCtrl.UpdateUserInfo) })

		// Shift

		// User
		api.GET("/user/shift/single/:shiftid", func(c *gin.Context) { fbCtrl.CheckAuth(c, shiftCtrl.GetUserSingle) })
		api.GET("/user/shift/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, shiftCtrl.GetUserAll) })
		api.POST("/user/shift/new", func(c *gin.Context) { fbCtrl.CheckAuth(c, shiftCtrl.Store) })

		// Admin
		api.GET("/shift/single/:shiftid", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, shiftCtrl.GetSingle) })
		api.GET("/shift/all", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, shiftCtrl.GetAll) })

		// Vouchers

		// User
		api.GET("/user/voucher/single/:voucherid", func(c *gin.Context) { fbCtrl.CheckAuth(c, vouchersCtrl.GetUserSingle) })
		api.GET("/user/voucher/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, vouchersCtrl.GetUserAll) })
		api.POST("/user/voucher/new", func(c *gin.Context) { fbCtrl.CheckAuth(c, vouchersCtrl.Store) })

		// Admin
		api.GET("/voucher/single/:voucherid", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, vouchersCtrl.GetSingle) })
		api.GET("/voucher/all", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, vouchersCtrl.GetAll) })

		// Deposits

		// User
		api.GET("/user/deposits/single/:depositid", func(c *gin.Context) { fbCtrl.CheckAuth(c, depositsCtrl.GetUserSingle) })
		api.GET("/user/deposits/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, depositsCtrl.GetUserAll) })
		api.POST("/user/deposits/new", func(c *gin.Context) { fbCtrl.CheckAuth(c, depositsCtrl.Store) })

		// Admin
		api.GET("/deposits/single/:depositid", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, depositsCtrl.GetSingle) })
		api.GET("/deposits/all", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, depositsCtrl.GetAll) })

		// Orders

		// User
		api.GET("/user/orders/single/:orderid", func(c *gin.Context) { fbCtrl.CheckAuth(c, ordersCtrl.GetUserSingle) })
		api.GET("/user/orders/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, ordersCtrl.GetUserAll) })
		api.POST("/user/orders/new", func(c *gin.Context) { fbCtrl.CheckAuth(c, ordersCtrl.Store) })

		// Admin
		api.GET("/orders/single/:orderid", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, ordersCtrl.GetSingle) })
		api.GET("/orders/all", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, ordersCtrl.GetAll) })

		// Cards

		// User
		api.GET("/user/cards/single/:cardcode", func(c *gin.Context) { fbCtrl.CheckAuth(c, cardsCtrl.GetUserSingle) })
		api.GET("/user/cards/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, cardsCtrl.GetUserAll) })

		// Admin
		api.GET("/cards/single/:cardcode", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, cardsCtrl.GetSingle) })
		api.GET("/cards/all", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, cardsCtrl.GetAll) })
		api.POST("/cards/new", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, cardsCtrl.Store) })

		// Coins Information

		// Admin
		api.GET("/coins", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, coinsCtrl.GetCoinsAvailability) })
		api.POST("/coins", func(c *gin.Context) { fbCtrl.CheckAuthAdmin(c, coinsCtrl.UpdateCoinsAvailability) })
	}
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})
}
