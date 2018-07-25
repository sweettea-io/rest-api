package model

import (
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
)

/*
  has_many --> Sessions
*/
type User struct {
  gorm.Model
  Uid        string    `gorm:"type:varchar(240);default:null;not null;unique;index:user_uid"`
  Email      string    `gorm:"type:varchar(240);default:null;not null;unique;index:user_email"`
  HashedPw   string    `gorm:"type:varchar(240);default:null"`
  Admin      bool      `gorm:"default:false"`
  Sessions   []Session
}

// Assign Uid to User before creation.
func (user *User) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", unique.NewUid())
  return nil
}