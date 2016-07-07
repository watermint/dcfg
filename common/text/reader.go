package text

import (
	"io/ioutil"
	"strings"
)

func SplitLines(rawText string) (lines []string) {
	// CRLF -> LF
	text := strings.Replace(rawText, "\r\n", "\n", -1)

	// CR -> LF
	text = strings.Replace(text, "\r", "\n", -1)

	return strings.Split(text, "\n")
}

func ReadLines(filePath string) (lines []string, err error) {
	seq, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []string{}, err
	}
	text := DecodeUnicodeSequence(seq)
	return SplitLines(text), nil
}

func ReadLinesIgnoreWhitespace(filePath string) (lines []string, err error) {
	rawLines, err := ReadLines(filePath)
	if err != nil {
		return []string{}, err
	}
	for _, rawLine := range rawLines {
		line := strings.TrimSpace(rawLine)
		if line != "" {
			lines = append(lines, strings.TrimSpace(line))
		}
	}
	return lines, nil
}
