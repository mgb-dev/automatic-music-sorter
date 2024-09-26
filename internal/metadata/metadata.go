package metadata

import (
	"slices"
	"strings"
)

type CriteriaType string

type Tags struct {
	title       string
	artist      string
	albumArtist string
	raw         map[string]string
}

type Metadata interface {
	Title() string
	Artist() string
	AlbumArtist() string
	Raw() map[string]interface{}
}

func (a *Tags) Title() string {
	return a.raw["title"]
}

func (a *Tags) AlbumArtist() string {
	return a.raw["albumArtist"]
}

func (a *Tags) Artist() string {
	return a.raw["artist"]
}

func (a *Tags) Raw() map[string]string {
	return a.raw
}

var criteriaList []CriteriaType = []CriteriaType{"albumartist", "artist", "title"}

func IsValidCriteria(str string) bool {
	return slices.Contains(criteriaList, CriteriaType(strings.ToLower(str)))
}
