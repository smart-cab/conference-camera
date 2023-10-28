package core

import (
	"conferencecam/ptz"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

type Api struct {
	App *gin.Engine
	Log *logrus.Logger
}

// Initializes gorm, gin, logrus and run http server
func Run() {
	api := Api{}

	// Init logrus
	api.Log = &logrus.Logger{
		Out:   os.Stdout,
		Level: logrus.DebugLevel,
		Formatter: &easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%\n",
		},
	}

	// Init .env config
	if err := godotenv.Load(); err != nil {
		api.Log.Fatal("failed load env config")
	}

	if err := ptz.Init(); err != nil {
		api.Log.Fatalf("failed start camera: %s", err.Error())
	}

	// Init gin engine
	gin.SetMode(gin.ReleaseMode)
	api.App = gin.Default()
	api.Routes()

	api.Log.Println("Server started!")
	// Run http server
	api.App.Run(fmt.Sprintf(
		"%s:%s",
		os.Getenv("IP"),
		os.Getenv("PORT"),
	))

}
