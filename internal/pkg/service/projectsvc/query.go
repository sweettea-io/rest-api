package projectsvc

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/model"
)

// All returns all Project records ordered by nsp.
func All() []model.Project {
  projects := []model.Project{}

  // Find Projects and eager-load ProjectConfig.
  app.DB.
    Preload("ProjectConfig").
    Order("nsp desc").
    Find(&projects)

  return projects
}

// FromNsp attempts to find a Project record for the given nsp.
// Will return an error if no record is found.
func FromNsp(nsp string, preloads ...string) (*model.Project, error) {
  var project model.Project
  result := app.DB

  for _, rel := range preloads {
    result = result.Preload(rel)
  }

  result = result.Where(&model.Project{Nsp: nsp}).Find(&project)

  // Return error if not found.
  if result.RecordNotFound() {
    return nil, fmt.Errorf("Project(nsp=%s) not found.\n", nsp)
  }

  return &project, nil
}

// FromID attempts to find a Project record by the provided id.
// Will return an error if no record is found.
func FromID(id uint) (*model.Project, error) {
  // Find Project by ID.
  var project model.Project
  result := app.DB.First(&project, id)

  // Return error if not found.
  if result.RecordNotFound() {
    return nil, fmt.Errorf("Project(ID=%v) not found.\n", id)
  }

  return &project, nil
}

func Filter(nsp string) []model.Project {
  if nsp == "" {
    return All()
  }

  var projects []model.Project
  project, _ := FromNsp(nsp)

  if project != nil {
    projects = append(projects, *project)
  }

  return projects
}
