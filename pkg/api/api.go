package api

import (
	"github.com/gin-gonic/gin"
	"go-proxypool/pkg/global"
	"net/http"
	"strconv"
)

type Server struct {
	port int
}

func NewServer(port int) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) Run() error {
	r := gin.Default()

	r.GET("/all-proxies", func(c *gin.Context) {
		allProxies, err := global.Storage.GetAll(c)
		if err != nil {
			global.Logger.Errorf("unable to retrieve all proxies: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve proxy IPs"})
		}
		c.JSON(http.StatusOK, allProxies)
	})

	return r.Run(":" + strconv.Itoa(s.port))
}
