package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	router := gin.Default()

	// corsの設定
	setCors(router)

  // ===== セッション（Cookie）設定 =====
	// 本番では必ず十分長いSECRETにする（32byte以上推奨）
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "CHANGE_ME_CHANGE_ME_CHANGE_ME_32BYTES"
	}
	store := cookie.NewStore([]byte(secret))

	// ローカルHTTPなら false、本番HTTPSなら true
	secureCookie := os.Getenv("COOKIE_SECURE") == "false"

	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7, // 7日
		HttpOnly: true,
		Secure:   secureCookie,
		SameSite: http.SameSiteLaxMode,
	})

	router.Use(sessions.Sessions("sid", store))


	router.POST("/signup", postSingup)
	router.POST("/login", postLogin)

	router.GET("/me", getMe)
	router.POST("/logout", postLogout)

	router.GET("/shop", getShopList)
  router.GET("/review/:shopId", getReview)

	auth := router.Group("/")
	auth.Use(requireLogin())
	{
		auth.PATCH("/name", patchName)
		auth.PATCH("/password", patchPassword)
		auth.POST("/review", postReview)
		// deleteはURL型が推奨：DELETE /review/:id
		auth.DELETE("/review/:id", deleteReviewByID)
	}

	return router
}

func setCors(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://ambitious-grass-0f2747200.2.azurestaticapps.net", "http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}