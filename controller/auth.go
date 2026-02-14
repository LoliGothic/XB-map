package controller

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/LoliGothic/XB-map/model"
)

func getMe(c *gin.Context) {
	sess := sessions.Default(c)
	v := sess.Get("userId")
	if v == nil {
		c.JSON(http.StatusUnauthorized, "not logged in")
		return
	}

	// 型変換（cookie storeの値がズレる可能性に備える）
	var userId int
	switch t := v.(type) {
	case int:
		userId = t
	case int64:
		userId = int(t)
	case float64:
		userId = int(t)
	default:
		c.JSON(http.StatusUnauthorized, "invalid session")
		return
	}

	user, err := model.GetUserByID(userId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "user not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.Id,
		"name":  user.Name,
		"email": user.Email,
	})
}

func postLogout(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Clear()
	sess.Options(sessions.Options{MaxAge: -1})
	_ = sess.Save()
	c.JSON(http.StatusOK, gin.H{"ok": true})
}