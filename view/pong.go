package view

import (
	"net/http"

	"github.com/NCNUCodeOJ/BackendClassManagement/models"

	"github.com/gin-gonic/gin"
)

func Pong(c *gin.Context) {
	if models.Ping() != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
