package main

// InSlice haystack slice 中に needle が含まれるか
func InSlice(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
