package web

import (
	"fmt"
	"net/http"
	"strings"
	"tinder-like-app/config"
	"tinder-like-app/container"
	"tinder-like-app/controller/web/handler"

	"github.com/gin-gonic/gin"
)

type HttpServer interface {
	Start() error
	GetHandler() (http.Handler, error)
}

type httpServer struct {
	config      config.Config
	engine      *gin.Engine
	controllers controllers
}

type controllers struct {
	user handler.User
}

func NewHttpServer(container *container.Container) *httpServer {
	gin.SetMode(gin.ReleaseMode)
	if strings.ToLower(container.Config().Service.LogLevel) == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	controllers := controllers{
		*handler.NewUser(container),
	}
	requestHandler := &httpServer{container.Config(), r, controllers}
	requestHandler.setupRouting()

	return requestHandler
}

func (h *httpServer) Start() error {
	return h.engine.Run(fmt.Sprintf("%s:%s", h.config.Service.Host, h.config.Service.Port))
}

func (h *httpServer) GetHandler() (http.Handler, error) {
	return h.engine, nil
}
