package utils

import (
	"strings"
)

func NormalizeDirName(str string) string {
	l := strings.ToLower(str)
	res := strings.ReplaceAll(l, " ", "-")
	return res
}
