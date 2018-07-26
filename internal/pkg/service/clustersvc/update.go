package clustersvc

import (
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
  "github.com/sweettea-io/rest-api/internal/pkg/util/str"
  "github.com/sweettea-io/rest-api/internal/pkg/util/cloud"
)

func SetName(cluster *model.Cluster, name string) {
  if name == "" {
    return
  }

  cluster.Name = name
  cluster.Slug = str.Slugify(name)
}

func SetCloud(cluster *model.Cluster, cloudName string) {
  if cloudName == "" || !cloud.IsValidCloud(cloudName) {
    return
  }

  cluster.Cloud = cloudName
}

func SetState(cluster *model.Cluster, state string) {
  if state == "" && !app.Config.OnTest() && !app.Config.OnLocal() {
    return
  }

  cluster.State = state
}