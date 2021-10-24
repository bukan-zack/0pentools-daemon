package router

import (
	"github.com/0pentools/daemon/domain"
	"github.com/0pentools/daemon/handler"
	"github.com/gin-gonic/gin"
)

func Configure(m *domain.Manager) *gin.Engine {
	r := gin.New()

	dh := handler.NewDomainHandler(m)

	r.Use(gin.Recovery())
	
	dg := r.Group("/domains")
	{
		dg.GET("/", dh.List)
		dg.POST("/", dh.Create)
		dg.GET("/start/:uuid", dh.Start)
	}

	return r
}