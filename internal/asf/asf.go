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
	asfObjGuidSize      = 16
	asfObjSize          = 8
	asfFileHeaderSize   = 8
	asfFileHeaderUnused = 6
)

var noMatchingGuidError = errors.New("No match found for given GUID")

type ByteSequence struct {
	ObjType asfObject
	Start   int
	End     int
}

func IsAsf(fileHeader []byte) bool {
	n := bytes.Compare(fileHeader, asfFileHeaderGuid)
	return n == 0
}

type AsfTags struct {
	raw map[string]string
}

func (this *AsfTags) Title() string {
	return this.raw["title"]
}

func (this *AsfTags) AlbumArtist() string {
	return this.raw["albumArtist"]
}

func (this *AsfTags) Artist() string {
	return this.raw["artist"]
}

func (this *AsfTags) Raw() *map[string]string {
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

type asfObject int

const (
	unknownObj asfObject = iota
	extendedContentDescriptionObj
	contentDescriptionObj
)

func parseContentDescription(bufPtr *[]byte, seq *ByteSequence) *AsfTags {
	buf := (*bufPtr)[seq.Start:seq.End]
	t := new(AsfTags)
	mT := make(map[string]string)
	// Descriptors count: 2 bytes
	for i := 2; i < len(buf); {
	}
}

func parseExtendedContentDescription(bufPtr *[]byte, seq *ByteSequence) *AsfTags {
	// Structure of data:
	// GUID: 16 bytes (already removed)
	// obj size: 8 bytes (already removed)
	buf := (*bufPtr)[seq.Start:seq.End]
	t := new(AsfTags)
	mT := make(map[string]string)
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
	t.raw = mT

	return t
}

func parseAsfObj(bufPtr *[]byte, seq *ByteSequence) *AsfTags {
	switch seq.ObjType {
	case extendedContentDescriptionObj:
		return parseExtendedContentDescription(bufPtr, seq)
	case contentDescriptionObj:
		return nil

	default:
		return nil
	}
}

// Receives
func findAsfObject(dPtr *[]byte, sPtr *[]ByteSequence, asfObjectType asfObject) error {
	// asf header objs have this form [16 bytes GUID][8 bytes obj size including header][data n bytes]
	guid := *getGuid(asfObjectType)

	data := *dPtr
	bSeq := *new(ByteSequence)

	for i := 0; i < len(data); {
		header := data[i:(i + asfObjGuidSize)]
		i += asfObjGuidSize
		size := int(binary.LittleEndian.Uint16(data[i:(i + asfObjSize - 1)]))
		i += asfObjSize
		if bytes.Compare(header, guid) != 0 {
			i = size
			continue
		}

		bSeq.ObjType = asfObjectType
		bSeq.Start = i
		bSeq.End = i + (size - asfObjGuidSize - asfObjSize)
		break

	}

	if bSeq.End == 0 {
		return noMatchingGuidError
	}

	return nil
}

func ReadAsf(rPtr *io.ReadSeeker) (*AsfTags, error) {
	r := *rPtr
	if _, err := r.Seek(asfObjGuidSize, io.SeekStart); err != nil {
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
