package util

import (
	"runtime"
	"strconv"
	"strings"
)

func FileWithLineNum() string {
	for i := 2; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok {
			split := "pkg" + strings.Split(file, "pkg")[1]

			return split + ":" + strconv.FormatInt(int64(line), 10)
		}
	}

	return ""
}
