package controller

import (
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	router := gin.Default()

	// corsの設定
	setCors(router)

	router.PATCH("/name", patchName)
	router.PATCH("/password", patchPassword)
	router.POST("/signup", postSingup)
	router.POST("/login", postLogin)
	router.POST("/check", postCheckLogin) //ログイン状態をチェックする
	router.POST("/review", postReview)
	router.GET("/shop", getShopList)
	router.GET("/review/:shopId", getReview)
	router.DELETE("/review", deleteReview)

	return router
}

func setCors(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://gentle-tree-0b58de200.2.azurestaticapps.net/", "http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}