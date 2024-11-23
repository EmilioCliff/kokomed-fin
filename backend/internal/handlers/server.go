package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) runSetup() {
	r := gin.Default()

	r.GET("/health-check", s.healthCheckHandler)

	s.router = r
}

func (s *Server) healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

func (s *Server) Stop() error {}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
