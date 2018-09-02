package testutil

import "github.com/sweettea-io/rest-api/internal/pkg/service/projectsvc"

func CreateProject(nsp string) RequestModifier {
  return func(req *Request) (*Request, error) {
    // Upsert project by namespace.
    _, _, err := projectsvc.UpsertByNsp(nsp)

    if err != nil {
      return nil, err
    }

    return req, nil
  }
}