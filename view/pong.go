package view

import (
	"NCNUOJBackend/ClassManagement/models"
	"net/http"

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
