package payload

import "net/http"

type Payload interface {
  Validate(req *http.Request) bool
}