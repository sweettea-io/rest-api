package err

import "github.com/gin-gonic/gin"

func Unauthorized() (int, *gin.H) {
  return 401, &gin.H{"ok": false, "code": 401, "error": "unauthorized"}
}

func Forbidden() (int, *gin.H) {
  return 403, &gin.H{"ok": false, "code": 403, "error": "forbidden"}
}

func UnknownError() (int, *gin.H) {
  return 500, &gin.H{"ok": false, "code": 500, "error": "unknown_error"}
}