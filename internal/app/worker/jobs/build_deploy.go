package jobs

import (
  "github.com/sweettea-io/work"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/util/stcluster"
  "fmt"
)

func (c *Context) BuildDeploy(job *work.Job) error {
  // Extract args from job.
  resourceID := job.ArgString("resourceID")
  projectID := uint(job.ArgInt64("projectID"))
  targetCluster := job.ArgString("targetCluster")

  if err := job.ArgError(); err != nil {
    app.Log.Errorln(err.Error())
    return err
  }

  // Validate target cluster.
  if targetCluster != stcluster.Train && targetCluster != stcluster.Api {
    err := fmt.Errorf("build deploy error: target cluster \"%s\" unsupported", targetCluster)
    app.Log.Errorln(err.Error())
    return err
  }

  return nil
}
