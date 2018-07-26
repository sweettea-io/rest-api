package model

import (
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
)

/*
  has_one -- Project
*/
type ProjectConfig struct {
  gorm.Model
}

func (pc *ProjectConfig) AsJSON() enc.JSON {
  return enc.JSON{}
}