package server

import (
	"clientapp/config"
	"clientapp/handler"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Config *config.Config
	Router *gin.Engine
}

func New() *Server {
	config, err := config.Load()
	if err != nil {
		log.Fatal("failed to start the server:", err)
	}

	server := &Server{
		Config: config,
	}

	if err := server.createRoutersAndSetRoutes(); err != nil {
		log.Fatal("failed to start the server:", err)
	}

	return server

}

func (s *Server) Run() {
	add := fmt.Sprintf("%s:%s", s.Config.Host, s.Config.Port)
	log.Println("[SERVER] listening on", add)

	s.Router.Run(add)
}

func (s *Server) createRoutersAndSetRoutes() error {
	handler := handler.New()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "route not found"})
	})

	router.POST("/login", handler.Login)

	s.Router = router
	return nil
}
