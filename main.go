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
		fbCtrl := controllers.FirebaseController{App: fbApp, UsersModel: usersModel}
		cardsCtrl := controllers.CardsController{Model: cardsModel}
		depositsCtrl := controllers.DepositsController{Model: depositsModel}
		ordersCtrl := controllers.OrdersController{Model: ordersModel}
		shiftCtrl := controllers.ShiftsController{Model: shiftsModel, Obol: obol, UserModel: usersModel}
		userCtrl := controllers.UsersController{Model: usersModel}
		vouchersCtrl := controllers.VouchersController{Model: vouchersModel}
		coinsCtrl := controllers.CoinsController{Model: coinsModel}

		// Coins Information
		// User
		api.GET("/coins", func(c *gin.Context) { fbCtrl.CheckAuth(c, coinsCtrl.GetCoinsAvailability, false) })
		// Admin
		api.POST("/coins", func(c *gin.Context) { fbCtrl.CheckAuth(c, coinsCtrl.UpdateCoinsAvailability, true) })

		// User Information
		// User
		api.GET("/user/info", func(c *gin.Context) { fbCtrl.CheckAuth(c, userCtrl.GetSelfInfo, false) })
		// Admin
		api.GET("/info/:uid", func(c *gin.Context) { fbCtrl.CheckAuth(c, userCtrl.GetUserInfo, true) })
		api.POST("/info", func(c *gin.Context) { fbCtrl.CheckAuth(c, userCtrl.UpdateUserInfo, true) })

		// Shift
		// User
		api.GET("/user/shift/single/:shiftid", func(c *gin.Context) { fbCtrl.CheckAuth(c, shiftCtrl.GetUserSingle, false) })
		api.GET("/user/shift/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, shiftCtrl.GetUserAll, false) })
		// TODO post body JWE decryption
		api.POST("/user/shift/new", func(c *gin.Context) { fbCtrl.CheckAuth(c, shiftCtrl.Store, false) })

		// Admin
		api.GET("/shift/single/:shiftid", func(c *gin.Context) { fbCtrl.CheckAuth(c, shiftCtrl.GetSingle, true) })
		api.GET("/shift/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, shiftCtrl.GetAll, true) })
		// TODO post body JWE decryption
		api.POST("/shift/update", func(c *gin.Context) { fbCtrl.CheckAuth(c, shiftCtrl.Update, true) })

		// Vouchers
		// User
		api.GET("/user/voucher/single/:voucherid", func(c *gin.Context) { fbCtrl.CheckAuth(c, vouchersCtrl.GetUserSingle, false) })
		api.GET("/user/voucher/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, vouchersCtrl.GetUserAll, false) })
		// TODO post body JWE decryption
		api.POST("/user/voucher/new", func(c *gin.Context) { fbCtrl.CheckAuth(c, vouchersCtrl.Store, false) })
		// Admin
		api.GET("/voucher/single/:voucherid", func(c *gin.Context) { fbCtrl.CheckAuth(c, vouchersCtrl.GetSingle, true) })
		api.GET("/voucher/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, vouchersCtrl.GetAll, true) })

		// Deposits
		// User
		api.GET("/user/deposits/single/:depositid", func(c *gin.Context) { fbCtrl.CheckAuth(c, depositsCtrl.GetUserSingle, false) })
		api.GET("/user/deposits/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, depositsCtrl.GetUserAll, false) })
		// TODO post body JWE decryption
		api.POST("/user/deposits/new", func(c *gin.Context) { fbCtrl.CheckAuth(c, depositsCtrl.Store, false) })
		// Admin
		api.GET("/deposits/single/:depositid", func(c *gin.Context) { fbCtrl.CheckAuth(c, depositsCtrl.GetSingle, true) })
		api.GET("/deposits/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, depositsCtrl.GetAll, true) })

		// Orders
		// User
		api.GET("/user/orders/single/:orderid", func(c *gin.Context) { fbCtrl.CheckAuth(c, ordersCtrl.GetUserSingle, false) })
		api.GET("/user/orders/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, ordersCtrl.GetUserAll, false) })
		// TODO post body JWE decryption
		api.POST("/user/orders/new", func(c *gin.Context) { fbCtrl.CheckAuth(c, ordersCtrl.Store, false) })
		// Admin
		api.GET("/orders/single/:orderid", func(c *gin.Context) { fbCtrl.CheckAuth(c, ordersCtrl.GetSingle, true) })
		api.GET("/orders/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, ordersCtrl.GetAll, true) })

		// Cards
		// User
		api.GET("/user/cards/single/:cardcode", func(c *gin.Context) { fbCtrl.CheckAuth(c, cardsCtrl.GetUserSingle, false) })
		api.GET("/user/cards/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, cardsCtrl.GetUserAll, false) })
		// Admin
		api.GET("/cards/single/:cardcode", func(c *gin.Context) { fbCtrl.CheckAuth(c, cardsCtrl.GetSingle, true) })
		api.GET("/cards/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, cardsCtrl.GetAll, true) })
		// TODO post body JWE decryption
		api.POST("/cards/new", func(c *gin.Context) { fbCtrl.CheckAuth(c, cardsCtrl.Store, true) })

	}
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})
}
