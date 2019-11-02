package main

import (
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	corsConf := cors.DefaultConfig()
	corsConf.AllowAllOrigins = true
	corsConf.AllowHeaders = []string{"token", "service", "content-type"}
	App.Use(cors.New(corsConf))
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

	// Init Database
	firestore, err := fbApp.Firestore(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	doc := firestore.Collection("polispay").Doc("hestia")

	// Init DB models
	shiftsModel := &models.ShiftModel{Firestore: doc, Collection: "shifts"}
	cardsModel := &models.CardsModel{Firestore: doc, Collection: "cards"}
	ordersModel := &models.OrdersModel{Firestore: doc, Collection: "orders"}
	depositsModel := &models.DepositsModel{Firestore: doc, Collection: "deposits"}
	vouchersModel := &models.VouchersModel{Firestore: doc, Collection: "vouchers"}
	usersModel := &models.UsersModel{Firestore: doc, Collection: "users"}
	coinsModel := &models.CoinsModel{Firestore: doc, Collection: "coins"}
	globalConfigModel := &models.GlobalConfigModel{Firestore: doc, Collection: "config"}
	exchangesModel := &models.ExchangesModel{Firestore: doc, Collection: "exchanges"}

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
	exchangesCtrl := controllers.ExchangesController{Model: exchangesModel}

	api := r.Group("/")
	{
		// Any
		api.GET("/user/coins", func(c *gin.Context) { fbCtrl.CheckAuth(c, coinsCtrl.GetCoinsAvailability, false) })
		api.GET("/user/config", func(c *gin.Context) { fbCtrl.CheckAuth(c, globalConfigCtrl.GetConfig, false) })
		api.GET("/user/info", func(c *gin.Context) { fbCtrl.CheckAuth(c, userCtrl.GetSingle, false) })
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
		// Admin
		api.POST("/coins", func(c *gin.Context) { fbCtrl.CheckAuth(c, coinsCtrl.UpdateCoinsAvailability, true) })
		api.POST("/config", func(c *gin.Context) { fbCtrl.CheckAuth(c, globalConfigCtrl.UpdateConfigData, true) })
		// TODO route
		api.GET("/balances/admin", func(c *gin.Context) { fbCtrl.CheckAuth(c, globalConfigCtrl.GetConfig, true) })

		// TODO pending routes
		//api.GET("/users/info/single/:uid", func(c *gin.Context) { fbCtrl.CheckAuth(c, userCtrl.GetSingle, true) })
		//.GET("/users/info/all", func(c *gin.Context) { fbCtrl.CheckAuth(c, userCtrl.GetAll, true) })
	}

	authUser := os.Getenv("HESTIA_AUTH_USERNAME")
	authPassword := os.Getenv("HESTIA_AUTH_PASSWORD")
	authApi := r.Group("/", gin.BasicAuth(gin.Accounts{
		authUser: authPassword,
	}))
	{
		// Routes available for another service
		// Tyche
		authApi.GET("/shift/single/:shiftid", shiftCtrl.GetSingleTyche)
		authApi.GET("/shift/all", shiftCtrl.GetAllTyche)
		authApi.POST("/shift", shiftCtrl.Store)

		// Ladon
		authApi.GET("/voucher/single/:voucherid", vouchersCtrl.GetSingleLadon)
		authApi.GET("/voucher/all", vouchersCtrl.GetAllLadon)
		authApi.POST("/voucher", vouchersCtrl.Store)

		// Adrestia
		authApi.GET("/adrestia/orders", exchangesCtrl.GetOrders)
		authApi.POST("/adrestia/new", exchangesCtrl.StoreOrder)
		authApi.PUT("/adrestia/new", exchangesCtrl.UpdateOrder)

		// TODO pending MS
		// Deposit Service
		//api.GET("/deposit/single/:depositid", depositsCtrl.GetSingle)
		//api.GET("/deposit/all", depositsCtrl.GetAll)
		//authApi.POST("/deposit", depositsCtrl.Store)

		// Order Service
		//authApi.GET("/order/single/:orderid", ordersCtrl.GetSingle)
		//authApi.GET("/order/all", ordersCtrl.GetAll)
		//authApi.POST("/order", ordersCtrl.Store)

		// Cards Service
		//api.GET("/card/single/:cardcode", cardsCtrl.GetSingle)
		//api.GET("/card/all", cardsCtrl.GetAll)
		//authApi.POST("/card", cardsCtrl.Store)

		// For all microservices
		api.GET("/coins", coinsCtrl.GetCoinsAvailabilityMicroService)
		api.GET("/config", globalConfigCtrl.GetConfigMicroservice)
		// TODO route
		api.GET("/balances/services", globalConfigCtrl.GetConfigMicroservice)
		authApi.POST("/validate/token", fbCtrl.CheckToken)
	}
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})
}
