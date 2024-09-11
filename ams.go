package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/dhowden/tag"
)

func ReadBytes(reader io.Reader, n int) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := reader.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

const (
	asfHeaderObjOffset = 30
)

func getASFHeaderObject() []byte {
	return []byte{
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
}

// ASF_Extended_Content_Description_Object
func getASFExtendedContentDescriptionObject() []byte {
	return []byte{
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
}

/* ByteSequence represents two pointers related to a []byte
 * StartOffset int: number of bytes from the ([]byte)[0] in which the desired sequence begins
 * Length: from StartOffset
 */
type ByteSequence struct {
	Startindex int
	Endindex   int
}

// Given a ASF GUID, finds the header, and returns a ByteSequence with its starting point and length
func FindASFObject(buf *[]byte, guid []byte) (*ByteSequence, error) {
	if len(guid) != 16 {
		return nil, errors.New("Provided GUID is not of correct size")
	}
	const header_len = 16
	const size_len = 8
	dref_buf := *buf
	byte_seq := new(ByteSequence)

	for i := 0; i < len(dref_buf); {

		header := dref_buf[i : i+header_len]
		size := int(binary.LittleEndian.Uint16(dref_buf[i+header_len : i+header_len+size_len]))
		if bytes.Compare(header, guid) != 0 {
			i = size
			continue
		}

		byte_seq.Startindex = i + header_len + size_len
		byte_seq.Endindex = byte_seq.Startindex + size - (header_len + size_len)
		break

	}

	if byte_seq.Endindex == 0 {
		return nil, errors.New("No match found for given GUID")
	}
	return byte_seq, nil
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

func ParseExtendedContentDescription(buf *[]byte) map[string]string {
	// Structure of data:
	// GUID: 16 bytes (already removed)
	// obj size: 8 bytes (already removed)
	dref_buf := *buf
	m := make(map[string]string)
	// Descriptors count: 2 bytes
	for i := 2; i < len(dref_buf); {

		descNameLen := toInt(dref_buf[i:i+2], "LE")
		i += 2
		descName := removeNullChar(string(dref_buf[i : i+descNameLen]))
		i += descNameLen + 2
		descValueLen := toInt(dref_buf[i:i+2], "LE")
		i += 2
		descValue := removeNullChar(string(dref_buf[i : i+descValueLen]))
		i += descValueLen
		key := strings.ReplaceAll(descName, "WM/", "")
		value := strings.ToLower(strings.ReplaceAll(descValue, " ", "-"))

		m[string(key)] = value

	}
	return m
}

// Removes first 31 bytes from the header object
func RemoveASFTopHeaderObject(buf []byte) []byte {
	const remove_from int = 30

	return buf[remove_from:]
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing argument: pos 1 - 'work_dir'")
		return
	}
	working_dir := os.Args[1]

	dir_entry, err := os.ReadDir(working_dir)
	if err != nil {
		log.Fatal("Directory reading error: \n", err)
	}
	fmt.Println("Working Directory :", working_dir)

	for _, file_entry := range dir_entry {
		if file_entry.IsDir() {
			continue
		}

		fmt.Println("File: ", file_entry.Name())
		file, err := os.Open(working_dir + file_entry.Name())
		defer file.Close()
		if err != nil {
			log.Fatal("File opening error:\n", err)
		}

		mdata, err := tag.ReadFrom(file)
		mdata.Artist()
		if err != nil && !errors.Is(err, tag.ErrNoTagsFound) {
			log.Fatal("tag reading error:\n", err)
		}
		buffer, err := ReadBytes(file, 24)
		if err != nil {
			log.Fatal(err)
		}

		top_level_header_obj := buffer[0:16]
		if bytes.Compare(getASFHeaderObject(), top_level_header_obj) != 0 {
			println("Program terminated")
			break
		}
		header_obj_size := int(binary.LittleEndian.Uint16(buffer[16:])) - asfHeaderObjOffset
		header_obj_data := make([]byte, header_obj_size)
		if _, err := file.ReadAt(header_obj_data, asfHeaderObjOffset); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Header Object size: %v bytes\n", len(header_obj_data))
		seq, err := FindASFObject(&header_obj_data, getASFExtendedContentDescriptionObject())
		if err != nil {
			log.Fatal(err)
		}
		asf_object_data := header_obj_data[seq.Startindex:seq.Endindex]

		contentDescriptorsMap := ParseExtendedContentDescription(&asf_object_data)
		fmt.Println("Artist: ", contentDescriptorsMap["AlbumArtist"])

	}
}
