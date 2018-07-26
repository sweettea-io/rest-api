package projectsvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

// All returns all Project records ordered by nsp
func All() []model.Project {
  projects := []model.Project{}
  app.DB.Order("nsp desc").Find(&projects)
  return projects
}

// UpsertByNsp upserts a Project by nsp.
// Return values are the project, if a new project was created, and any error that occurred, respectively.
func UpsertByNsp(nsp string) (*model.Project, bool, error) {
  isNewProj := false

  // Find Project by nsp.
  var project model.Project
  result := app.DB.Where(&model.Project{Nsp: nsp}).Find(&project)

  // If not found, create the Project/ProjectConfig.
  if result.RecordNotFound() {
    isNewProj = true

    // Create db transaction.
    tx := app.DB.Begin()

    // Create ProjectConfig.
    pc := model.ProjectConfig{}

    // Rollback if needed.
    if err := tx.Create(&pc).Error; err != nil {
      tx.Rollback()
      return nil, isNewProj, fmt.Errorf("error creating ProjectConfig: %s", err.Error())
    }

    // Create Project.
    project = model.Project{
      Nsp: nsp,
      ProjectConfig: &pc,
    }

    // Rollback if needed.
    if err := tx.Create(&project).Error; err != nil {
      tx.Rollback()
      return nil, isNewProj, fmt.Errorf("error creating Project: %s", err.Error())
    }

    // Commit writes.
    if err := tx.Commit().Error; err != nil {
      return nil, isNewProj, fmt.Errorf("error creating Project/ProjectConfig pair: %s", err.Error())
    }
  }

  return &project, isNewProj, nil
}