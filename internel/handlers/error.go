package handlers

import "github.com/gin-gonic/gin"

func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
