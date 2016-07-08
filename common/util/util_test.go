package util

import (
	"testing"
)

func TestUtilContainsString(t *testing.T) {
	stringHaystack := []string{"Simplify", "The", "Way", "People", "Work", "Together"}
	if !ContainsString(stringHaystack, "Simplify") {
		t.Errorf("Simplify not found")
	}
	if ContainsString(stringHaystack, "Complexity") {
		t.Errorf("Complexity found!")
	}
}
