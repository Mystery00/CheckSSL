package handler

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	"CheckSSL/utils"
)

func TestHandler(ctx *gin.Context) {
	utils.LogInfof(ctx, "receive test request")
	jsonBytes, _ := json.Marshal(struct {
		Result  bool
		Message string
	}{true, "Check Success"})
	ctx.Writer.Write(jsonBytes)
}
