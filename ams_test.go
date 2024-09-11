package main

import (
	"bytes"
	"testing"
)

func TestReadBytes(t *testing.T) {
	const to_read int = 5
	data := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21}
	want := data[0:to_read]
	reader := bytes.NewReader(want)
	got, _ := ReadBytes(reader, to_read)

	if bytes.Compare(want, got) != 0 {
		t.Errorf("got % X want % X", got, want)
	}
}

func TestRemoveHeaderObject(t *testing.T) {
	// hello world in hex
	want := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0x57, 0x6f, 0x72, 0x6c, 0x64, 0x21}
	extra_data := []byte{
		0x08,
		0x90,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
	}
	mock_data := append(getASFHeaderObject(), extra_data...)
	mock_data = append(mock_data, want...)
	got := RemoveASFTopHeaderObject(mock_data)

	if bytes.Compare(want, got) != 0 {
		t.Errorf("got % X want % X", got, want)
	}
}

var extended_content_obj = []byte{
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

func getMockASFData() *[]byte {
	extra_header := []byte{
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
	extra_size := []byte{0x7c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	extra_data := make([]byte, 100)
	ext_con_obj_size := []byte{0x25, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	ext_con_obj_data := []byte{
		0x48,
		0x65,
		0x6c,
		0x6c,
		0x6f,
		0x2c,
		0x20,
		0x57,
		0x6f,
		0x72,
		0x6c,
		0x64,
		0x21,
	}
	buf := append(extra_header, extra_size...)
	buf = append(buf, extra_data...)
	buf = append(buf, extended_content_obj...)
	buf = append(buf, ext_con_obj_size...)
	buf = append(buf, ext_con_obj_data...)
	return &buf
}

func TestFindASFObject(t *testing.T) {
	mockData := getMockASFData()
	drefData := *mockData
	want := "Hello, World!"
	byteSeq, err := FindASFObject(mockData, extended_content_obj)
	if err != nil {
		t.Errorf("FindASFObject error: %v", err)
		return
	}

	got := drefData[byteSeq.Startindex:byteSeq.Endindex]

	if string(got) != want {
		t.Errorf("got %s want %s", got, want)
	}
}
