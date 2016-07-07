package file

import "testing"

func TestFileExist(t *testing.T) {
	notFound := "/noexistent"
	found := "file_test.go"

	if FileExist(notFound) {
		t.Errorf("%s found!", notFound)
	}
	if !FileExist(found) {
		t.Error("%s not found", found)
	}
}

func TestFileExistAndReadable(t *testing.T) {
	notFound := "/noexistent"
	directory := "."
	found := "file_test.go"

	if FileExistAndReadable(notFound) {
		t.Errorf("%s found!", notFound)
	}
	if FileExistAndReadable(directory) {
		t.Errorf("%s should not marked as readable", directory)
	}
	if !FileExistAndReadable(found) {
		t.Error("%s not found", found)
	}
}

func TestIsDirectory(t *testing.T) {
	directory := "."
	file := "file_test.go"

	if IsDirectory(file) {
		t.Errorf("%s should not marked as directory", file)
	}
	if !IsDirectory(directory) {
		t.Errorf("%s should be marked as directory", directory)
	}
}
