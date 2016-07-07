package text

import (
	"bytes"
	"unicode/utf16"
)

var (
	bomUtf8    = []byte{0xef, 0xbb, 0xbf}
	bomUtf16BE = []byte{0xfe, 0xff}
	bomUtf16LE = []byte{0xff, 0xfe}
	//bomUtf32BE = []byte{0x00, 0x00, 0xfe, 0xff}
	//bomUtf32LE = []byte{0xff, 0xfe, 0x00, 0x00}
)

func toUint16Seq(seq []byte, bigEndian bool) []uint16 {
	utf16Seq := make([]uint16, len(seq)/2)
	var x, y int
	if bigEndian {
		x = 0
		y = 1
	} else {
		x = 1
		y = 0
	}
	for i := range utf16Seq {
		high := uint16(seq[2 * i + x]) << 8
		low := uint16(seq[2 * i + y])
		utf16Seq[i] = high | low
	}
	return utf16Seq
}

func DecodeUnicodeSequence(seq []byte) string {
	if bytes.HasPrefix(seq, bomUtf8) {
		return string(bytes.TrimPrefix(seq, bomUtf8))
	}
	if bytes.HasPrefix(seq, bomUtf16BE) {
		seqWithoutBom := bytes.TrimPrefix(seq, bomUtf16BE)
		return string(utf16.Decode(toUint16Seq(seqWithoutBom, true)))

	}
	if bytes.HasPrefix(seq, bomUtf16LE) {
		seqWithoutBom := bytes.TrimPrefix(seq, bomUtf16LE)
		return string(utf16.Decode(toUint16Seq(seqWithoutBom, false)))
	}
	return string(seq)
}
