package companysvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app"
  "errors"
)

func FindBySlug(slug string) (*model.Company, error) {
  var company model.Company
  result := app.DB.Where(&model.Company{Slug: slug}).First(&company)

  // Return error if company doesn't exist.
  if result.RecordNotFound() {
    return nil, errors.New("")
  }

  return &company, nil
}
