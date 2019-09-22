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

		// Init DB models
		shiftsModel := &models.ShiftModel{Db: db, Collection: "shifts"}
		cardsModel := &models.CardsModel{Db: db, Collection: "cards"}
		ordersModel := &models.OrdersModel{Db: db, Collection: "orders"}
		depositsModel := &models.DepositsModel{Db: db, Collection: "deposits"}
		vouchersModel := &models.VouchersModel{Db: db, Collection: "vouchers"}
		usersModel := &models.UsersModel{Db: db, Collection: "users"}
		coinsModel := &models.CoinsModel{Db: db, Collection: "coins"}
		globalConfigModel := &models.GlobalConfigModel{Db: db, Collection: "config"}

		// Init Controllers
		fbCtrl := controllers.FirebaseController{App: fbApp, UsersModel: usersModel}
		cardsCtrl := controllers.CardsController{Model: cardsModel, UserModel: usersModel}
		depositsCtrl := controllers.DepositsController{Model: depositsModel, UserModel: usersModel}
		ordersCtrl := controllers.OrdersController{Model: ordersModel, UserModel: usersModel}
		shiftCtrl := controllers.ShiftsController{Model: shiftsModel, UserModel: usersModel}
		userCtrl := controllers.UsersController{Model: usersModel}
		vouchersCtrl := controllers.VouchersController{Model: vouchersModel, UserModel: usersModel}
		coinsCtrl := controllers.CoinsController{Model: coinsModel}
		globalConfigCtrl := controllers.GlobalConfigController{Model: globalConfigModel}

		// Routes available for users
		api.GET("/coins", func(c *gin.Context) { fbCtrl.CheckAuth(c, coinsCtrl.GetCoinsAvailability, false) })
		api.GET("/config", func(c *gin.Context) { fbCtrl.CheckAuth(c, globalConfigCtrl.GetConfig, false) })
		api.GET("/user/info/single/:uid", func(c *gin.Context) { fbCtrl.CheckAuth(c, userCtrl.GetSingle, false) })
		api.GET("/user/shift/single/:shiftid", func(c *gin.Context) { fbCtrl.CheckAuth(c, shiftCtrl.GetSingle, false) })
		api.GET("/user/shift/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, shiftCtrl.GetAll, false) })
		api.GET("/user/voucher/single/:voucherid", func(c *gin.Context) { fbCtrl.CheckAuth(c, vouchersCtrl.GetSingle, false) })
		api.GET("/user/voucher/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, vouchersCtrl.GetAll, false) })
		api.GET("/user/deposit/single/:depositid", func(c *gin.Context) { fbCtrl.CheckAuth(c, depositsCtrl.GetSingle, false) })
		api.GET("/user/deposit/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, depositsCtrl.GetAll, false) })
		api.GET("/user/card/single/:cardcode", func(c *gin.Context) { fbCtrl.CheckAuth(c, cardsCtrl.GetSingle, false) })
		api.GET("/user/card/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, cardsCtrl.GetAll, false) })
		api.GET("/user/order/single/:orderid", func(c *gin.Context) { fbCtrl.CheckAuth(c, ordersCtrl.GetSingle, false) })
		api.GET("/user/order/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, ordersCtrl.GetAll, false) })

		// Routes available for admin users
		api.POST("/coins", func(c *gin.Context) { fbCtrl.CheckAuth(c, coinsCtrl.UpdateCoinsAvailability, true) })
		api.POST("/config", func(c *gin.Context) { fbCtrl.CheckAuth(c, globalConfigCtrl.UpdateConfigData, true) })

		api.GET("/deposit/single/:depositid", func(c *gin.Context) { fbCtrl.CheckAuth(c, depositsCtrl.GetSingle, true) })
		api.GET("/deposit/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, depositsCtrl.GetAll, true) })

		api.GET("/shift/single/:shiftid", func(c *gin.Context) { fbCtrl.CheckAuth(c, shiftCtrl.GetSingle, true) })
		api.GET("/shift/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, shiftCtrl.GetAll, true) })

		api.GET("/user/info/single/:uid", func(c *gin.Context) { fbCtrl.CheckAuth(c, userCtrl.GetSingle, true) })
		api.POST("/user/info/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, userCtrl.GetAll, true) })

		api.GET("/voucher/single/:voucherid", func(c *gin.Context) { fbCtrl.CheckAuth(c, vouchersCtrl.GetSingle, true) })
		api.GET("/voucher/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, vouchersCtrl.GetAll, true) })

		api.GET("/card/single/:cardcode", func(c *gin.Context) { fbCtrl.CheckAuth(c, cardsCtrl.GetSingle, true) })
		api.GET("/card/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, cardsCtrl.GetAll, true) })

		api.GET("/order/single/:orderid", func(c *gin.Context) { fbCtrl.CheckAuth(c, ordersCtrl.GetSingle, true) })
		api.GET("/order/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, ordersCtrl.GetAll, true) })

		// Routes available for another service

		// Shift Service
		api.POST("/shift", shiftCtrl.Store)

		// Voucher Service
		api.POST("/voucher", vouchersCtrl.Store)

		// Deposit Service
		api.POST("/deposit", depositsCtrl.Store)

		// Order Service
		api.POST("/order", ordersCtrl.Store)

		// Cards Service
		api.POST("/card", cardsCtrl.Store)
	}
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})
}
