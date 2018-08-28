package testutil

import (
  "fmt"
  "github.com/sweettea-io/rest-api/internal/app"
  "github.com/sweettea-io/rest-api/internal/pkg/service/sessionsvc"
  "github.com/sweettea-io/rest-api/internal/pkg/service/usersvc"
  "github.com/sweettea-io/rest-api/internal/pkg/util/unique"
)

func defaultTestUserEmail() string {
  return fmt.Sprintf("testuser-%s@sweettea.io", unique.NewUid()[:8])
}

func defaultTestUserHashedPw() string {
  return ""
}

func defaultTestUserAdmin() bool {
  return false
}

func AuthReqWithNewUser(req *Request) (*Request, error) {
  // Create new test user.
  user, err := usersvc.Create(
    defaultTestUserEmail(),
    defaultTestUserHashedPw(),
    defaultTestUserAdmin(),
  )

  if err != nil {
    return nil, err
  }

  // Create new session for user.
  session, err := sessionsvc.Create(user)
  if err != nil {
    return nil, err
  }

  // Add session token to auth request header.
  req.SetHeader(app.Config.AuthHeaderName, session.Token)

  // Return updated request.
  return req, nil
}

