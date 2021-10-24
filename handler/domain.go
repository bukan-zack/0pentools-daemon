package handler

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/0pentools/daemon/domain"
)

type DomainHandler struct {
	manager *domain.Manager
}

type CreateDomainForm struct {
	Memory  uint `form:"memory" binding:"required"`
	Cpu  uint `form:"cpu" binding:"required"`
	Disk uint `form:"disk" binding:"required"`
}

func NewDomainHandler(m *domain.Manager) *DomainHandler {
	return &DomainHandler{
		manager: m,
	}
}

func (h *DomainHandler) List(c *gin.Context) {
	ds, err := h.manager.GetDomains()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get domains",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"domains": ds,
	})
}

func (h *DomainHandler) Create(c *gin.Context) {
	form := CreateDomainForm{}
	if err := c.ShouldBind(&form); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid form body",
		})
		return
	}

	d, err := h.manager.CreateDomain(uint64(form.Cpu), uint64(form.Memory))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to create domain",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"uuid":    d,
	})
}

func (h *DomainHandler) Start(c *gin.Context) {
	uuid := c.Param("uuid")
	if uuid == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "uuid not defined in request",
		})
		return
	}

	if err := h.manager.StartDomain(uuid); err != nil {
		if err.Error() == fmt.Sprintf(`Domain not found: no domain with matching name '%v'`, uuid) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "domain not found",
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to start domain",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}