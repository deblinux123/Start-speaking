package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func home(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "Server is running"})
}
