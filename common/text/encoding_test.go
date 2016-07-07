package text

import (
	"encoding/binary"
	"testing"
	"unicode/utf16"
)

func TestDecodeUnicodeSequence(t *testing.T) {
	text := "Test/テスト/試験"
	textRune := []rune(text)
	textUtf8 := []byte(text)
	textUtf16 := utf16.Encode(textRune)

	// create UTF16 Little Endian byte sequence
	textUtf16LeSeq := make([]byte, 0, len(textUtf16)*2+2)
	textUtf16LeSeq = append(textUtf16LeSeq, bomUtf16LE...)
	for _, x := range textUtf16 {
		s := make([]byte, 2)
		binary.LittleEndian.PutUint16(s, x)
		textUtf16LeSeq = append(textUtf16LeSeq, s...)
	}

	// create UTF16 Big Endian byte sequence
	textUtf16BeSeq := make([]byte, 0, len(textUtf16)*2+2)
	textUtf16BeSeq = append(textUtf16BeSeq, bomUtf16BE...)
	for _, x := range textUtf16 {
		s := make([]byte, 2)
		binary.BigEndian.PutUint16(s, x)
		textUtf16BeSeq = append(textUtf16BeSeq, s...)
	}

	if text != DecodeUnicodeSequence(textUtf8) {
		t.Error("Failed for UTF-8")
	}
	if s := DecodeUnicodeSequence(textUtf16LeSeq); text != s {
		t.Errorf("Failed for UTF-16LE [%s] <=> [%s]; %x", text, s, textUtf16LeSeq)
	}
	if s := DecodeUnicodeSequence(textUtf16BeSeq); text != s {
		t.Errorf("Failed for UTF-16BE [%s] <=> [%s]; %x", text, s, textUtf16BeSeq)
	}
}
