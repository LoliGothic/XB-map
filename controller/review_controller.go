package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/LoliGothic/XB-map/model"
)

type NewReviewRequest struct {
	ShopId      int    `json:"shopId"`
	Explanation string `json:"explanation"`
}

type DeleteReview struct {
	Id int
	ShopId int
}

func getReview(c *gin.Context) {
  shopId, _ := strconv.Atoi(c.Param("shopId"))

  reviews, err := model.ReviewList(shopId)
  if err != nil {
    c.JSON(500, err.Error())
    return
  }

  c.JSON(200, reviews)
}

func postReview(c *gin.Context) {
	// ① セッションから userId
	sess := sessions.Default(c)
	v := sess.Get("userId")
	if v == nil {
		c.JSON(http.StatusUnauthorized, "not logged in")
		return
	}

	// 型変換
	userId := 0
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

	// ② body から shopId / explanation だけ受け取る
	var req NewReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}

	// ③ modelへ（userIdはサーバが決める）
	review, err := model.AddReview(userId, req.ShopId, req.Explanation)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, review)
}


func deleteReviewByID(c *gin.Context) {
	// ① セッションから userId を取得（ログイン必須）
	sess := sessions.Default(c)
	v := sess.Get("userId")
	if v == nil {
		c.JSON(http.StatusUnauthorized, "not logged in")
		return
	}

	// 型変換（保存方式で型がズレることがある）
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

	// ② パラメータから reviewId
	idStr := c.Param("id")
	reviewId, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid review id")
		return
	}

	// ③ “投稿者本人だけ削除”を model 側で強制
	reviews, err := model.DeleteReviewByID(reviewId, userId)
	if err != nil {
		// err の種類で 403/404 などに分けてもいい
		c.JSON(http.StatusForbidden, err.Error())
		return
	}

	// ④ 更新後のレビュー一覧返す（今の仕様に合わせる）
	c.JSON(http.StatusOK, reviews)
}