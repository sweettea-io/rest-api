package model

import (
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
)

/*
  User <-- belongs_to
*/
type Session struct {
  gorm.Model
  User       User
  UserID     uint   `gorm:"default:null;not null;index:session_user_id"`
  Token      string `gorm:"type:varchar(240);default:null;not null;unique;index:session_token"`
}

// Assign newly minted secret to Session Token before creation.
func (session *Session) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Token", unique.FreshSecret())
  return nil
}