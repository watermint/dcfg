package domain

type Entity interface {
	Id() string
}

func ContainsString(haystack []string, needle string) bool {
	for _, x := range haystack {
		if needle == x {
			return true
		}
	}
	return false
}

func ContainsEntity(haystack []Entity, needle Entity) bool {
	for _, x := range haystack {
		if needle.Id() == x.Id() {
			return true
		}
	}
	return false
}

func UniqueEntity(haystack []Entity) (unique []Entity) {
	entities := make(map[string]Entity)
	for _, e := range haystack {
		entities[e.Id()] = e
	}
	for _, e := range entities {
		unique = append(unique, e)
	}
	return
}
