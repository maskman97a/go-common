package config

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-common/common"
	"go-common/config/environment"
	"go-common/config/http_server"
	"go-common/config/middleware"
	"go-common/config/server"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type App interface {
	Run()

	Init()

	AddController()
}

type app struct {
	env environment.Environment

	gin *gin.Engine

	gormDB *gorm.DB

	db *sql.DB

	httpServer server.HttpServer
}

func NewApp() App {
	return &app{}
}

func (app *app) Run() {

	app.httpServer.Start()
	defer app.httpServer.Stop()

	// Listen for OS signals to perform a graceful shutdown
	log.Println("listening signals...")
	c := make(chan os.Signal, 1)
	signal.Notify(
		c,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	<-c
	log.Println("graceful shutdown...")
}

func (app *app) Init() {
	app.env = environment.ConfigAppEnv()

	app.gin = gin.New()
	gin.SetMode(gin.ReleaseMode)

	app.gin.Use(gin.Recovery())

	providerConfig := common.ProvideDatabaseConfig()

	app.db, _ = common.NewSQLDB(providerConfig)

	app.httpServer = server.NewHttpServer(
		app.gin,
		http_server.HttpServerConfig{
			Port:        app.env.ServerPort,
			ContextPath: app.env.ContextPath,
		},
	)

	app.AddController()

}

func (app *app) AddController() {
	app.gin.Use(middleware.MiddleWare)
	contextGroup := app.gin.Group(app.httpServer.GetContextPath())
	logrus.Info("Context path: " + app.httpServer.GetContextPath())
	baseController := common.NewBaseController()
	baseController.InitRouter(contextGroup)

}
