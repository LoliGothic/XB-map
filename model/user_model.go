package model

import (
	"errors"
	"strings"
	"unicode/utf8"
	"github.com/LoliGothic/XB-map/crypto"
	"github.com/google/uuid"
)

type User struct {
	Id int `gorm:"primaryKey"`
	Uuid string
	Name string
	Password string
	Email string
	Reviews []Review
}

func Signup(name string, password string, checkPassword string, email string) (*User, error) {
	if utf8.RuneCountInString(name) <= 0 || utf8.RuneCountInString(name) > 20 {
		err := errors.New("ユーザー名は20文字以内で入力してください")
		return nil, err
	}
	if utf8.RuneCountInString(password) <= 5 {
		err := errors.New("パスワードは6文字以上で入力してください")
		return nil, err
	}
	if utf8.RuneCountInString(checkPassword) <= 5 {
		err := errors.New("パスワードは6文字以上で入力してください")
		return nil, err
	}
	if strings.Contains(email, "@") == false{
		err := errors.New("正しいメールアドレスを入力してください")
		return nil, err
	}

	user := User{}

	if password != checkPassword {
		err := errors.New("パスワードが一致しません")
		return nil, err
	}

	// SELECT * FROM user WHERE email = 'email' ORDER BY id LIMIT 1;をしている
	db.Where("email = ?", email).First(&user)
	if user.Id != 0 {
		err := errors.New("同一のemailがすでに登録されています")
		return nil, err
	}

	// パスワードを暗号化
	hash, err := crypto.PasswordEncrypt(password)
	if err != nil {
		err := errors.New("パスワード暗号化中にエラーが発生しました")
		return nil, err
	}

	// uuid作成
	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	
	user = User{Uuid: uuid, Name: name, Password: hash, Email: email}
	db.Create(&user)

	return &user, nil
}

func Login(password string, email string) (*User, error) {
	if utf8.RuneCountInString(password) <= 5 {
		err := errors.New("パスワードは6文字以上で入力してください")
		return nil, err
	}
	if strings.Contains(email, "@") == false{
		err := errors.New("正しいメールアドレスを入力してください")
		return nil, err
	}

	user := User{}
	
	// SELECT * FROM user WHERE email = 'email' ORDER BY id LIMIT 1;
	db.Where("email = ?", email).First(&user)
	if user.Id == 0 {
		err := errors.New("一致するユーザが存在しません")
		return nil , err
	}

	// パスワードが一致するか確認
	err := crypto.CompareHashAndPassword(user.Password, password)
	if err != nil {
		err := errors.New("パスワードが一致しません")
		return nil, err
	}

	return &user, nil
}

func ChangeName(email string, newName string) (error) {
	if utf8.RuneCountInString(newName) <= 0 || utf8.RuneCountInString(newName) > 20 {
		err := errors.New("ユーザー名は20文字以内で入力してください")
		return err
	}

	db.Model(&User{}).Where("email = ?", email).Update("name", newName)

	return nil
}

func ChangePassword(email string, currentPassword string, newPassword string, checkNewPassword string) (*User, error) {
	if utf8.RuneCountInString(currentPassword) <= 5 {
		err := errors.New("パスワードは6文字以上で入力してください")
		return nil, err
	}
	if utf8.RuneCountInString(newPassword) <= 5 {
		err := errors.New("パスワードは6文字以上で入力してください")
		return nil, err
	}
	if utf8.RuneCountInString(checkNewPassword) <= 5 {
		err := errors.New("パスワードは6文字以上で入力してください")
		return nil, err
	}

	user := User{}

	if newPassword != checkNewPassword {
		err := errors.New("新しいパスワードが一致しません")
		return nil, err
	}

	// SELECT * FROM user WHERE email = 'email' ORDER BY id LIMIT 1;
	db.Where("email = ?", email).First(&user)
	if user.Id == 0 {
		err := errors.New("一致するユーザが存在しません")
		return nil, err
	}

	// パスワードが一致するか確認
	err := crypto.CompareHashAndPassword(user.Password, currentPassword)
	if err != nil {
		err := errors.New("現在のパスワードが一致しません")
		return nil, err
	}

	// パスワードを暗号化
	hash, err := crypto.PasswordEncrypt(newPassword)
	if err != nil {
		err := errors.New("パスワード暗号化中にエラーが発生しました")
		return nil, err
	}

	db.Model(&user).Update("password", hash)

	return &user, nil
}

func GetUserByID(userId int) (User, error) {
  var user User
  if err := db.First(&user, userId).Error; err != nil {
    return User{}, err
  }
  return user, nil
}