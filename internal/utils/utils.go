package utils

import (
	"fmt"
	"slices"
	"strings"
)

func NormalizeDirName(str string) string {
	l := strings.ToLower(str)
	res := strings.ReplaceAll(l, " ", "-")
	res = strings.ReplaceAll(res, "/", "")
	return res
}

// Helper function that runs printf conditionally and returns the same boolean
// It's intended to be used in a if assignment
func ConditionalPrintf(bool bool, str string, any ...any) bool {
	if bool {
		// FIXME: unhandled "Printf" error case
		fmt.Printf(str, any...)
	}
	return bool
}

type DirList struct {
	list []string
}

func NewDirList() *DirList {
	return &DirList{list: make([]string, 0)}
}

func (ds *DirList) Exists(dir string) bool {
	if len(ds.list) < 1 {
		return false
	}
	return slices.Contains(ds.list, dir)
}

func (ds *DirList) Add(newDir string) {
	dsRef := ds.list
	if ds.Exists(newDir) {
		return
	}
	ds.list = append(dsRef, newDir)
	return
}
