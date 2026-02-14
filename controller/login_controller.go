package controller

import (
  "net/http"

  "github.com/gin-contrib/sessions"
  "github.com/gin-gonic/gin"
  "github.com/LoliGothic/XB-map/model"
)

type PreviousUser struct {
	Password string
	Email string
}

type NewUser struct {
	Name string
	Password string
	CheckPassword string
	Email string
}

func postLogin(c *gin.Context) {
  var previousUser PreviousUser
  if err := c.ShouldBindJSON(&previousUser); err != nil {
    c.JSON(http.StatusBadRequest, "invalid request")
    return
  }

  user, err := model.Login(previousUser.Password, previousUser.Email)
  if err != nil {
    c.JSON(http.StatusUnauthorized, "invalid email or password")
    return
  }

  sess := sessions.Default(c)
  sess.Set("userId", user.Id) // user.Id はあなたのUser構造体に合わせる
  if err := sess.Save(); err != nil {
    c.JSON(http.StatusInternalServerError, "failed to create session")
    return
  }

  c.JSON(http.StatusOK, gin.H{
    "id": user.Id,
    "name": user.Name,
    "email": user.Email,
  })
}

func postSingup(c *gin.Context) {
	var newUser NewUser //NewUser型の変数を定義
	c.BindJSON(&newUser) //受け取ったJSONをnewUserに代入
	user, err := model.Signup(newUser.Name, newUser.Password, newUser.CheckPassword, newUser.Email)
	if err == nil {
		c.JSON(200, user)
	} else {
		c.JSON(400, err.Error())
	}
}