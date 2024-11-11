package utils

import (
	"slices"
	"strings"
)

func NormalizeDirName(str string) string {
	l := strings.ToLower(str)
	res := strings.ReplaceAll(l, " ", "-")
	return res
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
