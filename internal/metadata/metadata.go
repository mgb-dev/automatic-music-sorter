package metadata

import (
	"io"
	"slices"
	"strings"

	"github.com/dhowden/tag"

	"github.com/mgb-dev/ams/internal/asf"
)

type CriteriaType string

type Metadata interface {
	Title() string
	Artist() string
	AlbumArtist() string
	Raw() *map[string]string
}

type Tags struct {
	raw map[string]string
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

func (a *Tags) Raw() *map[string]string {
	return &(a.raw)
}

var criteriaList []CriteriaType = []CriteriaType{"albumartist", "artist", "title"}

func IsValidCriteria(str string) bool {
	return slices.Contains(criteriaList, CriteriaType(strings.ToLower(str)))
}

// Converts tag.Metadata -> ams.metadata.Tags
func convert(mT *tag.Metadata) *Tags {
	result := map[string]string{}
	result["artist"] = (*mT).Artist()
	result["albumartist"] = (*mT).AlbumArtist()
	result["title"] = (*mT).Title()
	return &Tags{raw: result}
}

func ReadTags(r io.ReadSeeker) (Metadata, error) {
	fileHeader, err := asf.ReadBytes(&r, asf.AsfObjGuidSize)
	if err != nil {
		return nil, err
	}

	if asf.IsAsf(fileHeader) {
		// asf format
		m, err := asf.ReadAsf(&r)
		if err != nil {
			return nil, err
		}
		return m, nil
	}

	// Move Reader pointer back to io.SeekStart
	negativeOffset := int64(asf.AsfObjGuidSize * -1)
	if _, err := r.Seek(negativeOffset, io.SeekCurrent); err != nil {
		return nil, err
	}
	mt, err := tag.ReadFrom(r)
	if err != nil {
		return nil, err
	}
	t := convert(&mt)

	// TODO: learn more about: X does not implement Y (... method has a pointer receiver)
	return t, nil
}
