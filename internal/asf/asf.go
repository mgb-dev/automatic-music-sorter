package asf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"strings"
)

// Advanced File System's Magic Word
var asfFileHeaderGuid asfGuid = []byte{
	0x30,
	0x26,
	0xB2,
	0x75,
	0x8E,
	0x66,
	0xCF,
	0x11,
	0xA6,
	0xD9,
	0x00,
	0xAA,
	0x00,
	0x62,
	0xCE,
	0x6C,
}

var asfContentDescriptionGuid asfGuid = []byte{
	0x33,
	0x26,
	0xb2,
	0x75,
	0x8e,
	0x66,
	0xcf,
	0x11,
	0xa6,
	0xd9,
	0x00,
	0xaa,
	0x00,
	0x62,
	0xce,
	0x6c,
}

var asfExtendedContentDescriptionGuid asfGuid = []byte{
	0x40,
	0xa4,
	0xd0,
	0xd2,
	0x07,
	0xe3,
	0xd2,
	0x11,
	0x97,
	0xf0,
	0x00,
	0xa0,
	0xc9,
	0x5e,
	0xa8,
	0x50,
}

type asfGuid []byte

func getGuid(o asfObject) *asfGuid {
	switch o {
	case extendedContentDescriptionObj:
		return &asfExtendedContentDescriptionGuid
	case contentDescriptionObj:
		return &asfContentDescriptionGuid
	default:
		return new(asfGuid)

	}
}

// GUID size  16 bytes(file header obj)
// size		  8 bytes(header obj size)
// unused   6 bytes
// = total 30 bytes
const (
	AsfObjGuidSize      = 16
	asfObjSize          = 8
	asfFileHeaderSize   = 8
	asfFileHeaderUnused = 6
)

var noMatchingGuidError = errors.New("No match found for given GUID")

type byteSequence struct {
	ObjType asfObject
	Start   int
	End     int
}

func IsAsf(fileHeader []byte) bool {
	n := bytes.Compare(fileHeader, asfFileHeaderGuid)
	return n == 0
}

type asfTags struct {
	raw map[string]string
}

func (this *asfTags) Title() string {
	return this.raw["title"]
}

func (this *asfTags) AlbumArtist() string {
	return this.raw["albumArtist"]
}

func (this *asfTags) Artist() string {
	return this.raw["artist"]
}

func (this *asfTags) Raw() *map[string]string {
	return &(this.raw)
}

func ReadBytes(rPtr *io.ReadSeeker, n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := (*rPtr).Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func toInt(buf []byte, mode string) int {
	if mode == "BE" {
		return int(binary.BigEndian.Uint16(buf))
	}

	return int(binary.LittleEndian.Uint16(buf))
}

func removeNullChar(s string) string {
	return strings.Replace(s, "\x00", "", -1)
}

func mergeMap(target *asfMap, m1 *asfMap, m2 *asfMap) {
	t := *target
	for k, v := range *m1 {
		t[k] = v
	}
	for k, v := range *m2 {
		t[k] = v
	}

	*target = t
}

type (
	asfObject int
	asfMap    map[string]string
)

const (
	unknownObj asfObject = iota
	extendedContentDescriptionObj
	contentDescriptionObj
)

func parseContentDescription(bufPtr *[]byte, seq *byteSequence) *asfMap {
	buf := (*bufPtr)[seq.Start:seq.End]
	i := 0

	tLen := toInt(buf[i:i+2], "LE")
	i += 2
	aLen := toInt(buf[i:i+2], "LE")
	// skip 8 bytes to read title
	i += 8
	title := removeNullChar(string(buf[i : i+tLen+1]))
	i += tLen
	author := removeNullChar(string(buf[i : i+aLen+1]))
	return &asfMap{"title": title, "artist": author}
}

func parseExtendedContentDescription(bufPtr *[]byte, seq *byteSequence) *asfMap {
	// Structure of data:
	// GUID: 16 bytes (already removed)
	// obj size: 8 bytes (already removed)
	buf := (*bufPtr)[seq.Start:seq.End]
	mT := asfMap{}
	// Descriptors count: 2 bytes
	for i := 2; i < len(buf); {

		descNameLen := toInt(buf[i:i+2], "LE")
		i += 2
		descName := removeNullChar(string(buf[i : i+descNameLen]))
		i += descNameLen + 2
		descValueLen := toInt(buf[i:i+2], "LE")
		i += 2
		descValue := removeNullChar(string(buf[i : i+descValueLen]))
		i += descValueLen
		key := strings.ToLower(strings.ReplaceAll(descName, "WM/", ""))
		value := strings.ToLower(strings.ReplaceAll(descValue, " ", "-"))

		mT[key] = value
	}
	return &mT
}

func parseAsfObj(bufPtr *[]byte, seq *byteSequence) *asfMap {
	switch seq.ObjType {
	case contentDescriptionObj:
		return parseContentDescription(bufPtr, seq)
	case extendedContentDescriptionObj:
		return parseExtendedContentDescription(bufPtr, seq)
	default:
		return nil
	}
}

// finds asfObjectType and returns a ByteSequence referencing said data
func findAsfObject(dPtr *[]byte, asfObjectType asfObject) (byteSequence, error) {
	// asf header objs have this form [16 bytes GUID][8 bytes obj size including header][data n bytes]
	guid := *getGuid(asfObjectType)

	data := *dPtr
	var bSeq byteSequence

	for i := 0; i < len(data); {
		header := data[i:(i + AsfObjGuidSize)]
		i += AsfObjGuidSize
		size := int(binary.LittleEndian.Uint16(data[i:(i + asfObjSize - 1)]))
		i += asfObjSize
		if bytes.Compare(header, guid) != 0 {
			i = size
			continue
		}

		bSeq.ObjType = asfObjectType
		bSeq.Start = i
		bSeq.End = i + (size - AsfObjGuidSize - asfObjSize)
		break

	}

	if bSeq.End == 0 {
		return bSeq, noMatchingGuidError
	}

	return bSeq, nil
}

func ReadAsf(rPtr *io.ReadSeeker) (*asfTags, error) {
	r := *rPtr
	if _, err := r.Seek(AsfObjGuidSize, io.SeekStart); err != nil {
		return nil, err
	}

	buf := make([]byte, asfFileHeaderSize)
	if _, err := r.Read(buf); err != nil {
		return nil, err
	}
	headerSize := int(binary.LittleEndian.Uint16(buf))
	if _, err := r.Seek(asfFileHeaderUnused, io.SeekCurrent); err != nil {
		return nil, err
	}
	dataSlice := make([]byte, headerSize)
	if _, err := r.Read(dataSlice); err != nil {
		return nil, err
	}

	obj := new(AsfTags)
	sByteSeq := new([]ByteSequence)

	err := findAsfObject(&dataSlice, extendedContentDescriptionObj)
	if errors.Is(err, noMatchingGuidError) {
		objSeq, err = findAsfObject(&dataSlice, contentDescriptionObj)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	obj = parseAsfObj(&dataSlice, objSeq)

	return obj, nil
}
