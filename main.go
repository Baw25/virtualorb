package main

import (
	"os"

	"github.com/Baw25/virtualorb/orb"

	"github.com/gin-gonic/gin"
)

const (
	V1Path = "/v1/virtualorb"
)

func SetupServerRoutes(router *gin.Engine) *gin.Engine {
	v1 := router.Group(V1Path)
	{
		v1.POST("/report", orb.PostStatusReport)
		v1.POST("/signup", orb.PostSignup)
		v1.POST("/signup/batch", orb.PostSignupBatch)
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
			})
		})
	}

	return router
}

func main() {
	router := gin.Default()
	SetupServerRoutes(router)
	port := os.Getenv("PORT")
	err := router.Run("localhost:" + port)

	if err != nil {
		panic(err)
	}
}
