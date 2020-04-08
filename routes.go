package gomato

import (
	"github.com/gin-gonic/gin"
	gcache "github.com/patrickmn/go-cache"
)

// Server represents the necessary components to run the server
type Server struct {
	Cache  *gcache.Cache
	Router *gin.Engine
}

// NewServer instantiates a Server object
func NewServer() *Server {
	cache := gcache.New(-1, -1) // cache with no expiration, no cleanup
	router := gin.Default()
	return &Server{Cache: cache, Router: router}
}

// BuildRoutes builds the routes to listen to slash commands
func (s *Server) BuildRoutes(router *gin.Engine) {
	timer := router.Group("/timer")
	{
		timer.POST("/start", s.startTimer)
		timer.POST("/pause", s.pauseTimer)
		timer.POST("/stop", s.stopTimer)
	}

	pomBreak := router.Group("/break")
	{
		pomBreak.POST("/start", s.startBreak)
		pomBreak.POST("/pause", s.pauseBreak)
		pomBreak.POST("/stop", s.stopBreak)
	}
}
