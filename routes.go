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
func (s *Server) BuildRoutes() {
	timer := s.Router.Group("/timer")
	{
		timer.POST("start", s.startTimer)
		timer.POST("pause", s.pauseTimer)
		timer.POST("resume", s.resumeTimer)
		timer.POST("stop", s.stopTimer)
	}
}
