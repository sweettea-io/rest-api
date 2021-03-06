package commitsvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/service/projectsvc"
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

// FromID attempts to find a Commit record by the provided id.
// Will return an error if no record is found.
func FromID(id uint) (*model.Commit, error) {
  // Find Commit by ID.
  var commit model.Commit
  result := app.DB.First(&commit, id)

  // Return error if not found.
  if result.RecordNotFound() {
    return nil, fmt.Errorf("Commit(ID=%v) not found.\n", id)
  }

  return &commit, nil
}

func FromShaOrLatest(sha string, project *model.Project) (*model.Commit, error) {
  // If sha exists, return the Commit for that sha.
  if sha != "" {
    return FromSha(sha)
  }

  // Fetch the latest sha from the remote project repo, and upsert the commit.
  host := projectsvc.GetHost(project)
  latestSha, err := host.LatestSha(project.Owner(), project.Repo())

  if err != nil {
    return nil, err
  }

  // Upsert Commit for fetched sha.
  commit, err := Upsert(project.ID, latestSha)
  if err != nil {
    return nil, err
  }

  return commit, nil
}