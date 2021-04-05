package tool

import "strings"

func SystemFromSymbol(symbol string) string {
	return strings.Split(symbol, "-")[0]
}
