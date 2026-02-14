package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func requireLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		sess := sessions.Default(c)
		if sess.Get("userId") == nil {
			c.JSON(http.StatusUnauthorized, "not logged in")
			c.Abort()
			return
		}
		c.Next()
	}
}