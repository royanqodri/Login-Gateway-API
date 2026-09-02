package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	docs "github.com/royanqodri/Login-Gateway-API/cmd/docs"
	repoTransaction "github.com/royanqodri/Login-Gateway-API/repository/transaction"
	"github.com/royanqodri/Login-Gateway-API/websocket"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"github.com/royanqodri/Login-Gateway-API/config"
	"github.com/royanqodri/Login-Gateway-API/controller"
	"github.com/royanqodri/Login-Gateway-API/database"
	"github.com/royanqodri/Login-Gateway-API/middleware"
	"github.com/royanqodri/Login-Gateway-API/repository"
	"github.com/royanqodri/Login-Gateway-API/service"
	"github.com/royanqodri/Login-Gateway-API/service/transaction"
	"github.com/royanqodri/Login-Gateway-API/util"
	"github.com/royanqodri/Login-Gateway-API/util/logging"
	"github.com/sirupsen/logrus"
)

var (
	// repository
	tUserRepo       repository.TUserRepository
	tUserModuleRepo repository.TUserModuleRepository

	tCustomerRepo repository.TCustomerRepository

	tNotificationsRepo repoTransaction.TNotificationsRepository

	// service
	sessionService service.SessionService
	loginService   service.LoginService
	tUserService   transaction.TUserService

	// controller
	// tUserModuleController controller.TUserModuleController
	loginController     controller.LoginController
	sessionController   controller.SessionController
	websocketController controller.WebSocketController
	tUserController     controller.TUserController
)

func main() {
	// init config
	initConfig()

	// init swagger
	initSwagger()

	// init location
	util.InitLocation()

	// init Redis
	// util.InitRedis()

	// init redis websocket
	websocket.InitWebsocketRedis()

	// init database
	initDatabase()

	// init log
	logging.InitLog()

	defer logging.CloseLog()

	// init repository
	initRepository()

	// init service
	initService()

	// init controller
	initController()

	// init router
	initRouter()

	// init FCM

	// init Telegram

}

func initConfig() {
	err := config.Init()
	if err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}
}

func initSwagger() {
	docs.SwaggerInfo.Title = "Login Gateway API"
	docs.SwaggerInfo.Description = "API Documentation for Login Service"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = config.Get().Swagger.SwaggerHost + "/gateway"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{config.Get().Swagger.SwaggerScheme}

	log.Println("Swagger initialized at", docs.SwaggerInfo.Host)
}

func initDatabase() {
	database.ConnectDatabase()
}

func initRepository() {
	tUserRepo = repository.NewTUserRepository()
	tUserModuleRepo = repository.NewTUserModuleRepository()

	tCustomerRepo = repository.NewTCustomerRepository()
}

func initService() {
	sessionService = service.NewSessionService(tUserModuleRepo, tUserRepo, util.RedisClient)
	loginService = service.NewLoginService(tUserRepo, tUserModuleRepo, tCustomerRepo)
	tUserService = transaction.NewTUserService(tUserRepo, tCustomerRepo)
}

func initController() {
	sessionController = controller.NewSessionController(sessionService)
	loginController = controller.NewLoginController(loginService)
	websocketController = *controller.NewWebSocketController()
	tUserController = controller.NewTUserController(tUserService)
}

// @title			Login Gateway API
// @version		1.0
// @description	API Documentation for Login Gateway Service
// @host			localhost:9001
// @BasePath		/
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func initRouter() {
	// set service mode
	if strings.ToLower(config.Get().ServiceMode) == "debug" {
		gin.SetMode(gin.DebugMode)
	} else if strings.ToLower(config.Get().ServiceMode) == "test" {
		gin.SetMode(gin.TestMode)
	} else if strings.ToLower(config.Get().ServiceMode) == "release" || strings.ToLower(config.Get().ServiceMode) == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// init router
	router := gin.New()
	router.SetTrustedProxies([]string{"0.0.0.0/0"})

	// public route - require no authentication
	router.Use(gin.Recovery(), util.CORSMiddleware())
	router.POST("/login", loginController.Login)
	router.POST("/login/google", loginController.LoginWithGoogle)
	router.POST("/login/facebook", loginController.LoginWithFacebook)

	router.POST("/user", tUserController.Post)
	router.GET("/user", tUserController.GetAll)
	// router.POST("/generate-token", loginController.OnboardGenerateToken)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Backend Server OK"})
	})
	//router.GET("/stats", middleware.StatsHandler)

	// swagger docs route
	if config.Get().Swagger.SwaggerEnable || gin.Mode() == gin.DebugMode {
		router.GET("/swagger/*any", middleware.SwaggerDynamicAuth(), ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// private route - requires authentication
	routerGroup := router.Group("/:customer_no")
	routerGroup.Use(middleware.AuthMiddleware())
	routerGroup.Use(middleware.LoggerMiddleware())
	//routerGroup.Use(middleware.CountMiddleware())

	routerGroup.GET("/session-data", sessionController.GetBySessionData)

	// websocket
	routerGroupWebSocket := router.Group("notification")
	routerGroupWebSocket.Use(middleware.LoggerMiddleware())
	routerGroupWebSocket.GET("/sub", websocketController.Subscribe)
	routerGroupWebSocket.GET("/pub", websocketController.Publish)
	// routerGroupWebSocket.POST("/sos", websocketController.Sos)
	routerGroupWebSocket.GET("/history", websocketController.GetNotifications)
	routerGroupWebSocket.POST("/post", websocketController.Post)

	// test notification
	routerGroupWebSocket.POST("/test", websocketController.Test)
	// routerGroupWebSocket.GET("/redis-ws", websocketController.RedisWebSocket)

	// graceful shutdown websocket
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
		log.Println("Shutting down server...")
		websocket.CleanupWebSocket()
		os.Exit(0)
	}()

	// init timeout
	readTimeout := time.Duration(config.Get().Service.ReadTimeout) * time.Second
	writeTimeout := time.Duration(config.Get().Service.WriteTimeout) * time.Second
	idleTimeout := time.Duration(config.Get().Service.IdleTimeout) * time.Second

	// init server
	server := &http.Server{
		Addr:         config.Get().Service.HttpAddress,
		Handler:      router,
		ReadTimeout:  readTimeout,  // Max time to read request body
		WriteTimeout: writeTimeout, // Max time to write response
		IdleTimeout:  idleTimeout,  // Keep-alive
	}

	logging.SendToDiscord(logging.FormatDiscordMessage(
		logging.INFO,
		logrus.Fields{
			"message": "Service is Running 🚀",
		}))

	// listen server
	err := server.ListenAndServe()
	util.PanicIfError(err)
}
