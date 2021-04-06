package tool

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

func SystemFromSymbol(symbol string) string {
	return strings.Split(symbol, "-")[0]
}

func GetContext(skip int) string {
	pc, file, no, ok := runtime.Caller(skip)
	if ok {
		fname := runtime.FuncForPC(pc).Name()
		var packageName string
		{
			firstDot := strings.Index(fname, ".")
			if firstDot == -1 {
				packageName = "<unknown>"
			} else {
				packageName = fname
			}
		}
		_, filename := filepath.Split(file)

		return fmt.Sprintf("%s:%d in %s", filename, no, packageName)
	}
	return "unable to determine caller"
}