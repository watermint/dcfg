package util

func ContainsString(haystack []string, needle string) bool {
	for _, x := range haystack {
		if needle == x {
			return true
		}
	}
	return false
}
