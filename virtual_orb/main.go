package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	V1Path = "/v1/virtualorb"
)

func SetupServerRoutes(router *gin.Engine) *gin.Engine {
	v1 := router.Group(V1Path)
	{
		v1.POST("/report", virtualorb.PostStatusReport)
		v1.POST("/report/batch", virtualorb.PostStatusReportBatch)
		v1.POST("/signup", virtualorb.PostSignup)
		v1.POST("/signup/batch", virtualorb.PostSignupBatch)
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
	// router.GET("/albums", getAlbums)
	SetupServerRoutes(router)
	port := viper.GetString("PORT")
	err := router.Run(port)
	if err != nil {
		panic(err)
	}
}
