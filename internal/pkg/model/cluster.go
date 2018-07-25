package model

import (
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/util/str"
)

/*
  has_many --> Deploys
*/
type Cluster struct {
  gorm.Model
  Name       string   `gorm:"type:varchar(360);default:null;not null"`
  Slug       string   `gorm:"type:varchar(360);default:null;not null;unique;index:cluster_slug"`
  Cloud      string   `gorm:"type:varchar(240);default:null;not null"`
  State      string   `gorm:"type:varchar(360);default:null"`
  Deploys    []Deploy
}

// Assign Slug to Cluster before creation.
func (cluster *Cluster) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Slug", str.Slugify(cluster.Name))
  return nil
}