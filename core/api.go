package core

import (
	"conferencecam/types"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Api struct {
	App *gin.Engine
}

// Initializes gorm, gin, logrus and run http server
func Run() {
	api := Api{}

	// Init .env config
	if err := godotenv.Load(); err != nil {
		types.Logger.Fatal("failed load env config")
	}

	// Init gin engine
	gin.SetMode(gin.ReleaseMode)
	api.App = gin.Default()
	api.Routes()

	types.Logger.Println("Server started!")
	// Run http server
	api.App.Run(fmt.Sprintf(
		"%s:%s",
		os.Getenv("IP"),
		os.Getenv("PORT"),
	))

}
