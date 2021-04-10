package login_handler

import (
	"github.com/gin-gonic/gin"
	"go-lab/app/go_web/go_gin/handler"
)

type LoginHandler struct {
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}
func (h *LoginHandler) UserLogin(ctx *gin.Context) {
	handler.ResponseOK(ctx)
}
