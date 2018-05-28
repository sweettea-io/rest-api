package e

type Error map[string]interface{}

func Forbidden() (int, *Error) {
  return 403, &Error{"ok": false, "code": 403, "error": "forbidden"}
}