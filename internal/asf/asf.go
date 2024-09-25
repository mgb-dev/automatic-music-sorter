package asf

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

