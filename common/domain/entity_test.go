package domain

import (
	"testing"
)

func TestCollectionsContainsString(t *testing.T) {
	stringHaystack := []string{"Simplify", "The", "Way", "People", "Work", "Together"}
	if !ContainsString(stringHaystack, "Simplify") {
		t.Errorf("Simplify not found")
	}
	if ContainsString(stringHaystack, "Complexity") {
		t.Errorf("Complexity found!")
	}
}

type CollectionsEntity struct {
	key   string
	value string
}

func (m CollectionsEntity) Id() string {
	return m.key
}

func TestCollectionsContainsEntity(t *testing.T) {
	haystack := []Entity{
		CollectionsEntity{key: "Simplify", value: "SIMPLIFY"},
		CollectionsEntity{key: "The", value: "THE"},
		CollectionsEntity{key: "Way", value: "WAY"},
		CollectionsEntity{key: "People", value: "PEOPLE"},
		CollectionsEntity{key: "Work", value: "WORK"},
		CollectionsEntity{key: "Together", value: "TOGETHER"},
	}
	needle1 := CollectionsEntity{key: "Simplify", value: "シンプリファイ"}
	needle2 := CollectionsEntity{key: "シンプリファイ", value: "Simplify"}

	if !ContainsEntity(haystack, needle1) {
		t.Errorf("needle1 not found!")
	}
	if ContainsEntity(haystack, needle2) {
		t.Errorf("needle2 found!")
	}
}

func TestCollectionsUniqueEntity(t *testing.T) {
	haystack := []Entity{
		CollectionsEntity{key: "Simplify", value: "SIMPLIFY"},
		CollectionsEntity{key: "The", value: "THE"},
		CollectionsEntity{key: "Way", value: "WAY1"},
		CollectionsEntity{key: "Way", value: "WAY2"},
		CollectionsEntity{key: "People", value: "PEOPLE"},
		CollectionsEntity{key: "People", value: "PEOPLE2"},
		CollectionsEntity{key: "Work", value: "WORK"},
		CollectionsEntity{key: "Together", value: "TOGETHER"},
	}
	unique := UniqueEntity(haystack)
	if len(unique) != 6 {
		t.Errorf("Invalid length: %d", len(unique))
	}
	expected := []string{"Simplify", "The", "Way", "People", "Work", "Together"}
	for _, x := range expected {
		if !ContainsEntity(unique, CollectionsEntity{key: x}) {
			t.Errorf("'%s' not found", x)
		}
	}
}
