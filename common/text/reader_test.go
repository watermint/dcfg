package text

import "testing"

func TestReadLinesIgnoreWhitespace(t *testing.T) {
	expectedLines := []string{
		"Test1/テスト/試験",
		"Test2/テスト/試験",
		"Test3/テスト/試験",
		"Test4/テスト/試験",
	}

	compareLines := func(x, y []string) bool {
		if len(x) != len(y) {
			return false
		}
		for i, q := range x {
			if q != y[i] {
				return false
			}
		}
		return true
	}

	compare := func(label, file string) {
		lines, err := ReadLinesIgnoreWhitespace(file)
		if err != nil {
			t.Errorf("%s; Read err: %s", label, err)
		}
		if !compareLines(expectedLines, lines) {
			t.Errorf("%s: Didn't match: %v", label, lines)
		}
	}

	compare("UTF-8", "reader_test_utf8.txt")
	compare("UTF-8 (with BOM)", "reader_test_utf8+bom.txt")
	compare("UTF-16LE (with BOM)", "reader_test_utf16le+bom.txt")
	compare("UTF-16BE (with BOM)", "reader_test_utf16be+bom.txt")
}
