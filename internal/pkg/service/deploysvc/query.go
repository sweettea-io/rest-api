package deploysvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

// FromID attempts to find a Deploy record by the provided id.
// Will return an error if no record is found.
func FromID(id uint) (*model.Deploy, error) {
  // Find Deploy by ID.
  var deploy model.Deploy
  result := app.DB.Preload("Commit").First(&deploy, id)

  // Return error if not found.
  if result.RecordNotFound() {
    return nil, fmt.Errorf("Deploy(ID=%v) not found.\n", id)
  }

  return &deploy, nil
}