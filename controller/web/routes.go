package web

import (
	"github.com/gin-gonic/gin"
)

func (h *httpServer) setupRouting() {
	router := h.engine

	router.GET("/ping", func(context *gin.Context) {
		context.String(200, "Ok")
	})

	// v1 API
	V1 := router.Group("/v1")
	{
		V1.POST("/signup", h.controllers.user.SignUp)
		V1.POST("/login", h.controllers.user.Login)
	}
}
