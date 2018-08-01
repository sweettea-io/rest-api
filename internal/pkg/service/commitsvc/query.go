package commitsvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app"
  "fmt"
)

// FromID attempts to find a Commit record by the provided sha.
// Will return an error if no record is found.
func FromSha(sha string) (*model.Commit, error) {
  var commit model.Commit
  result := app.DB.Where(&model.Commit{Sha: sha}).Find(&commit)

  if result.RecordNotFound() {
    return nil, fmt.Errorf("Commit(sha=%s) not found.\n", sha)
  }

  return &commit, nil
}