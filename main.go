package main

import (
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"flag"
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

var polisPayDatabase string

func init() {
	_ = godotenv.Load()
}

func main() {
	// Read input flag
	localRun := flag.Bool("local", false, "set this flag to run hestia with testing data")
	flag.Parse()

	// If flag was set, change the polispay database to use testing data.
	if *localRun {
		polisPayDatabase = "hestia_test"
	} else {
		polisPayDatabase = "hestia"
	}

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
	doc := firestore.Collection("polispay").Doc(polisPayDatabase)
	bitcouDoc := firestore.Collection("bitcou")
	bitcouTestDoc := firestore.Collection("bitcou_test")
	bitcouConfDoc := firestore.Collection("bitcou_filters")

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
	AdrestiaModel := models.NewAdrestiaModel(*doc)
	balancesModel := &models.BalancesModel{Firestore: doc, Collection: "balances"}
	bitcouModel := &models.BitcouModel{Firestore: bitcouDoc, FirestoreTest: bitcouTestDoc}
	bitcouConfModel := &models.BitcouConfModel{Firestore: bitcouConfDoc}

	// Init Controllers
	fbCtrl := controllers.FirebaseController{App: fbApp, UsersModel: usersModel}
	cardsCtrl := controllers.CardsController{Model: cardsModel, UserModel: usersModel}
	depositsCtrl := controllers.DepositsController{Model: depositsModel, UserModel: usersModel}
	ordersCtrl := controllers.OrdersController{Model: ordersModel, UserModel: usersModel}
	shiftCtrl := controllers.ShiftsController{Model: shiftsModel, UserModel: usersModel}
	userCtrl := controllers.UsersController{Model: usersModel}
	AdrestiaCtrl := controllers.AdrestiaController{Model: &AdrestiaModel}

	vouchersCtrl := controllers.VouchersController{
		Model:           vouchersModel,
		UserModel:       usersModel,
		BitcouModel:     bitcouModel,
		BitcouConfModel: bitcouConfModel,
		CachedVouchers: controllers.VouchersCache{
			Vouchers:               make(map[string]controllers.CachedVouchersData),
			CachedCountries:        []string{},
			CachedCountriesUpdated: 0,
		}}
	coinsCtrl := controllers.CoinsController{Model: coinsModel, BalancesModel: balancesModel}
	globalConfigCtrl := controllers.GlobalConfigController{Model: globalConfigModel}
	exchangesCtrl := controllers.ExchangesController{Model: exchangesModel}
	statsCtrl := controllers.StatsController{ShiftModel: shiftsModel, VouchersModel: vouchersModel, DepositsModel: depositsModel, OrdersModel: ordersModel}

	// API Groups the endpoints that require Firebase token authentication.
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

		// Vouchers list
		api.GET("/user/voucher/list", func(c *gin.Context) { fbCtrl.CheckAuth(c, vouchersCtrl.GetAvailableCountries, false) })
		api.GET("/user/voucher/list/:country", func(c *gin.Context) { fbCtrl.CheckAuth(c, vouchersCtrl.GetVouchers, false) })
		// Voucher routes for development environment
		api.GET("/user/voucher/dev/list", func(c *gin.Context) { fbCtrl.CheckAuth(c, vouchersCtrl.GetTestAvailableCountries, false) })
		api.GET("/user/voucher/dev/list/:country", func(c *gin.Context) { fbCtrl.CheckAuth(c, vouchersCtrl.GetTestVouchers, false) })

		// Stats routes
		// Total Stats
		api.GET("/user/stats/shift", func(c *gin.Context) { fbCtrl.CheckAuth(c, statsCtrl.GetShiftStats, true) })
		api.GET("/user/stats/vouchers", func(c *gin.Context) { fbCtrl.CheckAuth(c, statsCtrl.GetVoucherStats, true) })

		// Admin
		api.POST("/coins", func(c *gin.Context) { fbCtrl.CheckAuth(c, coinsCtrl.UpdateCoinsAvailability, true) })
		api.POST("/config", func(c *gin.Context) { fbCtrl.CheckAuth(c, globalConfigCtrl.UpdateConfigData, true) })
		api.GET("/balances", func(c *gin.Context) { fbCtrl.CheckAuth(c, coinsCtrl.GetCoinBalances, true) })
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
		authApi.GET("/voucher/all_by_timestamp", vouchersCtrl.GetVouchersByTimestampLadon)

		// Adrestia
		authApi.GET("/adrestia/deposits", AdrestiaCtrl.GetDeposits)
		authApi.GET("/adrestia/withdrawals", AdrestiaCtrl.GetWithdrawals)
		authApi.GET("/adrestia/orders", AdrestiaCtrl.GetBalancerOrders)
		authApi.GET("/adrestia/balancer", AdrestiaCtrl.GetBalancers)
		authApi.POST("/adrestia/new/deposit", AdrestiaCtrl.StoreDeposit)
		authApi.POST("/adrestia/new/withdrawal", AdrestiaCtrl.StoreWithdrawal)
		authApi.POST("/adrestia/new/order", AdrestiaCtrl.StoreBalancerOrder)
		authApi.POST("/adrestia/new/balancer", AdrestiaCtrl.StoreBalancer)
		authApi.PUT("/adrestia/update/deposit", AdrestiaCtrl.UpdateDeposit)
		authApi.PUT("/adrestia/update/withdrawal", AdrestiaCtrl.UpdateWithdrawal)
		authApi.PUT("/adrestia/update/order", AdrestiaCtrl.UpdateBalancerOrder)
		authApi.PUT("/adrestia/update/balancer", AdrestiaCtrl.UpdateBalancer)

		// For all microservices
		api.GET("/exchange", exchangesCtrl.GetExchange)
		api.GET("/exchanges", exchangesCtrl.GetExchanges)
		api.POST("/exchanges/update", exchangesCtrl.UpdateExchange)
		api.GET("/coins", coinsCtrl.GetCoinsAvailabilityMicroService)
		api.GET("/config", globalConfigCtrl.GetConfigMicroservice)
		authApi.POST("/validate/token", fbCtrl.CheckToken)
	}
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})

	authAdminUser := os.Getenv("HESTIA_ADMIN_USER")
	authAdminPassword := os.Getenv("HESTIA_ADMIN_PASSWORD")
	authAdminApi := r.Group("/", gin.BasicAuth(gin.Accounts{
		authAdminUser: authAdminPassword,
	}))
	{
		authAdminApi.POST("/voucher/filter", vouchersCtrl.AddFilters)
		authAdminApi.DELETE("/voucher/filter", vouchersCtrl.AddFilters)
		authAdminApi.GET("/voucher/filter/:dev", vouchersCtrl.AddFilters)
	}
	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})
}
