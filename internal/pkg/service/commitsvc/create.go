package commitsvc

import (
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/app"
  "fmt"
  "github.com/sweettea-io/rest-api/internal/pkg/service/projectsvc"
)

func Upsert(projectID uint, sha string) (*model.Commit, error) {
  var commit model.Commit

  if err := app.DB.Where(&model.Commit{ProjectID: projectID, Sha: sha}).FirstOrCreate(&commit).Error; err != nil {
    return nil, fmt.Errorf("error upserting Commit: %s", err.Error())
  }

  return &commit, nil
}

func FetchAndUpsertFromSha(project *model.Project, sha string) (*model.Commit, error) {
  var err error
  host := projectsvc.GetHost(project)

  if sha == "latest" {
    sha, err = host.LatestSha(project.Owner(), project.Repo())
  } else {
    err = host.EnsureCommitExists(project.Owner(), project.Repo(), sha)
  }

  if err != nil {
    return nil, err
  }

  // Upsert Commit for fetched sha.
  commit, err := Upsert(project.ID, sha)
  if err != nil {
    return nil, err
  }

  return commit, nil
}