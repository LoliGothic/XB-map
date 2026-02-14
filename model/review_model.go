package model

import (
	"errors"
	"time"
	"unicode/utf8"
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	UserId int
	ShopId int
	Explanation string
}

type Result struct {
	Id int
	ShopId int
	UserId int
	CreatedAt time.Time
	Explanation string
	Name string
}

func ReviewList(shopId int) (*[]Result, error) {
	result := []Result{}

	err := db.Model(&Review{}).
		Select("reviews.id, reviews.shop_id, reviews.user_id, reviews.created_at, reviews.explanation, users.name").
		Joins("inner join users on users.id = reviews.user_id").
		Where("reviews.shop_id = ?", shopId).
		Order("reviews.created_at desc").
		Scan(&result).Error

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func AddReview(userId int, shopId int, explanation string) (*[]Result, error) {
	if utf8.RuneCountInString(explanation) <= 0 || utf8.RuneCountInString(explanation) > 100 {
		return nil, errors.New("口コミは100文字以内で入力してください")
	}

	review := Review{UserId: userId, ShopId: shopId, Explanation: explanation}

	// 新しい口コミをデータベースに追加（エラー見る）
	if err := db.Create(&review).Error; err != nil {
		return nil, err
	}

	result := []Result{}

	// password/emailは返さない。user_idは返す（削除ボタン判定に必要）
	if err := db.Model(&Review{}).
		Select("reviews.id, reviews.shop_id, reviews.user_id, reviews.created_at, reviews.explanation, users.name").
		Joins("inner join users on users.id = reviews.user_id").
		Where("reviews.shop_id = ?", shopId).
		Order("reviews.created_at desc").
		Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
)

func DeleteReviewByID(reviewId int, userId int) (*[]Result, error) {
	// ① 対象レビュー取得（shop_id と user_id を確定）
	var review Review
	if err := db.Select("id, shop_id, user_id").First(&review, reviewId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	// ② 投稿者本人かチェック
	if review.UserId != userId {
		return nil, ErrForbidden
	}

	// ③ 削除
	if err := db.Delete(&Review{}, reviewId).Error; err != nil {
		return nil, err
	}

	// ④ shopIdの口コミを全て返す（password/emailは返さない）
	result := []Result{}
	if err := db.Model(&Review{}).
		Select("reviews.id, reviews.shop_id, reviews.created_at, reviews.explanation, users.name, users.id as user_id").
		Joins("inner join users on users.id = reviews.user_id").
		Where("reviews.shop_id = ?", review.ShopId).
		Order("reviews.created_at desc").
		Scan(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}