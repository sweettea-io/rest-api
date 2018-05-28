package err

type E map[string]interface{}

//func Unauthorized() (int, *gin.H) {
//  return 401, &gin.H{"ok": false, "code": 401, "error": "unauthorized"}
//}

func Forbidden() (int, *E) {
  return 403, &E{"ok": false, "code": 403, "error": "forbidden"}
}
//
//func UnknownError() (int, *gin.H) {
//  return 500, &gin.H{"ok": false, "code": 500, "error": "unknown_error"}
//}