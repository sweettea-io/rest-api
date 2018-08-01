package buildable

import "github.com/sweettea-io/rest-api/internal/pkg/model"

type Buildable interface {
  GetCommit() *model.Commit
  GetUid() string
}