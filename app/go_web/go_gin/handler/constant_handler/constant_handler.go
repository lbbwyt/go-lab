package constant_handler

import (
	"github.com/gin-gonic/gin"
	"go-lab/app/go_web/go_gin/handler"
	"net/http"
)

// @Summary 获取所有常量
// @Description 获取所有常量
// @Tags constant
// @Accept  json
// @Produce  json
// @Success 200 {string} string "map[string]interface{}"
// @Router /api/all_constants [get]
func GetAllConstants(ctx *gin.Context) {
	constantsMap := make(map[string]interface{})

	handler.ResponseJSON(ctx, http.StatusOK, constantsMap)
}
