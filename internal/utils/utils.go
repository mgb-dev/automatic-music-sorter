package utils

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

func NormalizeDirName(str string) string {
	l := strings.ToLower(str)
	res := strings.ReplaceAll(l, " ", "-")
	res = strings.ReplaceAll(res, "/", "")
	return res
}

var invalidStringPath error = errors.New("Invalid string path")

// Expands `str` into a path string
func ExpandPath(str string) (string, error) {
	hp := strings.HasPrefix

	switch true {
	case hp(str, "./"):
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		path := path.Join(cwd, str[2:])
		return path, nil
	case hp(str, "~"):
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return homeDir, nil
	case hp(str, "~/"):
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path := path.Join(homeDir, str[2:])
		return path, nil
	case filepath.IsAbs(str):
		return str, nil
	case !(string(str[0]) == "/"):
		return filepath.Abs(str)
	default:
		return "", invalidStringPath
	}
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
