package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"CheckSSL/config"
	"CheckSSL/data"
	"CheckSSL/handler"
	"CheckSSL/utils"
)

var port = ":3000"

func init() {
	if !utils.FileExist(config.EnvConfig.ConfigFile) {
		log.Errorf("配置文件[%s]不存在", config.EnvConfig.ConfigFile)
		os.Exit(1)
	}
	gin.DisableConsoleColor()
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(func(context *gin.Context) {
		utils.Trace(context)
		context.Writer.Header().Set("Content-Type", "application/json")
	})
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%v] [%s] [%3d] [%13v] | %-7s %#v\n%s",
			param.TimeStamp.Format("2006-01-02 15:04:05.000"),
			param.Keys[utils.TraceId],
			param.StatusCode,
			param.Latency,
			param.Method,
			param.Path,
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())
	router.NoRoute(func(context *gin.Context) {
		jsonBytes, _ := json.Marshal(data.ErrorResponse{Code: http.StatusNotFound, Message: http.StatusText(http.StatusNotFound)})
		context.Writer.Write(jsonBytes)
	})

	router.GET("/test", handler.TestHandler)
	router.GET("/check", handler.CheckHandler)

	log.Infof("Server is running at %s", port)
	router.Run(port)
}
