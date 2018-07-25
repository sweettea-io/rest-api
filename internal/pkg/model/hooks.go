package model

import (
  "github.com/jinzhu/gorm"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
  "github.com/sweettea-io/rest-api/internal/pkg/util/str"
)

// WithUid should be added to each model with a string Uid
// column that needs to be auto-populated before creation.
type WithUid struct {}

// Assign Uid to model before creation.
func (w *WithUid) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", unique.NewUid())
  return nil
}

// Assign newly minted secret to Session token before creation.
func (session *Session) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Token", unique.FreshSecret())
  return nil
}

// Assign Uid & Slug to Model before creation.
func (model *Model) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", unique.NewUid())
  scope.SetColumn("Slug", str.Slugify(model.Name))
  return nil
}

// Assign Version to ModelVersion before creation.
func (mv *ModelVersion) BeforeCreate(scope *gorm.Scope) error {
  // TODO: Shorten version to 7 chars
  scope.SetColumn("Version", unique.NewUid())
  return nil
}

// Assign Uid, ClientID, & ClientSecret to Deploy before creation.
func (cluster *Deploy) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", unique.NewUid())
  scope.SetColumn("ClientID", unique.NewUid())
  scope.SetColumn("ClientSecret", unique.FreshSecret())
  return nil
}

// Assign Uid & Slug to Cluster before creation.
func (cluster *Cluster) BeforeCreate(scope *gorm.Scope) error {
  scope.SetColumn("Uid", unique.NewUid())
  scope.SetColumn("Slug", str.Slugify(cluster.Name))
  return nil
}