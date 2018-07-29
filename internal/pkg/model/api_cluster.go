package model

import (
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/util/enc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/str"
)

/*
  has_many --> Deploys
*/
type ApiCluster struct {
  gorm.Model
  Name       string   `gorm:"type:varchar(360);default:null;not null"`
  Slug       string   `gorm:"type:varchar(360);default:null;not null;unique;index:api_cluster_slug"`
  Cloud      string   `gorm:"type:varchar(240);default:null;not null"`
  State      string   `gorm:"type:varchar(360);default:null"`
  Deploys    []Deploy
}

// Assign Slug to Cluster before creation.
func (apiCluster *ApiCluster) BeforeSave(scope *gorm.Scope) error {
  scope.SetColumn("Slug", str.Slugify(apiCluster.Name))
  return nil
}

func (apiCluster *ApiCluster) AsJSON() enc.JSON {
  return enc.JSON{
    "name": apiCluster.Name,
    "slug": apiCluster.Slug,
    "cloud": apiCluster.Cloud,
    "state": apiCluster.State,
  }
}