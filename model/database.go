package model

import (
  "fmt"
  "os"

  "gorm.io/driver/mysql"
  "gorm.io/gorm"
)

var db *gorm.DB

func init() {
  dsn := os.Getenv("DB_DSN")
  if dsn == "" {
    panic("DB_DSN is not set")
  }

  var err error
  db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
  if err != nil {
    panic(fmt.Sprintf("DB connection failed: %v", err))
  }

  fmt.Println("データベース接続成功")

  if err := db.AutoMigrate(&User{}, &Shop{}, &Review{}); err != nil {
    panic(fmt.Sprintf("AutoMigrate failed: %v", err))
  }
}