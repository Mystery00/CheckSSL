package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

	"CheckSSL/config"
	"CheckSSL/data"
	"CheckSSL/utils"
)

func CheckHandler(ctx *gin.Context) {
	domainListJson := gjson.GetBytes(utils.ReadFile(config.EnvConfig.ConfigFile), "list").Array()
	var domainList []data.Domain
	for _, domainName := range domainListJson {
		domain := data.Domain{Name: domainName.String()}
		domain.DoCheckLocal(ctx)
		domainList = append(domainList, domain)
	}
	jsonBytes, err := json.Marshal(domainList)
	if err != nil {
		panic(err)
	}
	ctx.Writer.Write(jsonBytes)
}
