package utils

import "fmt"

func Assert(err error, messages ...string) {
  if err != nil {
    msg := ""

    // Only use the first message (hack for Go optional param)
    if len(messages) > 0 {
      msg = fmt.Sprintf("%s: ", messages[0])
    }

    panic(fmt.Errorf("%s%s", msg, err))
  }
}