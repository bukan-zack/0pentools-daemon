package web

import (
	"fmt"
	"net/http"

	"github.com/0penTools/panel/domain"
	"github.com/gin-gonic/gin"
)

type CreateDomainForm struct {
	Ram  uint `form:"ram" binding:"required"`
	Cpu  uint `form:"cpu" binding:"required"`
	Disk uint `form:"disk" binding:"required"`
}

// TODO: improve error messages

func ListDomain(ctx *gin.Context) {
	list, err := domain.GetAllDomains()
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Something went wrong",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"success": true,
		"domains": list,
	})
}

func CreateDomain(ctx *gin.Context) {
	form := CreateDomainForm{}
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Invalid form body",
		})
		return
	}
	dom, err := domain.Create(uint64(form.Cpu), uint64(form.Ram))
	if err != nil {
		ctx.JSON(500, gin.H{
			"success": false,
			"message": "Something went wrong",
		})
		fmt.Println(err)
		return
	}
	ctx.JSON(204, gin.H{
		"success": true,
		"uuid":    dom,
	})
}

func StartDomain(ctx *gin.Context) {
	err := domain.Start(ctx.PostForm("uuid"))
	if err != nil {
		if err.Error() == fmt.Sprintf(`Domain not found: no domain with matching name '%v'`, ctx.PostForm("uuid")) {
			ctx.JSON(404, gin.H{
				"message": "vm not found",
			})
			return
		}
		ctx.JSON(500, gin.H{
			"message": "Something went wrong",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "started",
	})
}
