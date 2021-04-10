package middleware

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

//前端获取用户所属团队的API接口不走鉴权，前端拿到用户团队信息后，
//将当前所选团队ID写入cookie（切换团队时刷新cookie）， 后续的接口请求均需携带cookie
func AuthMiddleWare(action string, handler gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sid, err := ctx.Cookie("sid")
		if err != nil || sid == "" {
			// cookie读取错误
			log.WithError(err).Error("get sid from cookie error")
			pleaseLoginResponse(ctx)
			return
		}
		//鉴权逻辑
		ctx.Set("current", "lbbwyt")
		handler(ctx)
		return
	}
}

func pleaseLoginResponse(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"msg": "请先登录账户",
	})
}
