package utils

// Remove deletes all occurrences of elem from the array
func Remove(arr []string, elem string) []string {
	a := make([]string, len(arr))
	copy(a, arr)
	for i := 0; i < len(a); i++ {
		if a[i] == elem {
			a[i] = a[len(a)-1]
			a[len(a)-1] = ""
			a = a[:len(a)-1]
		}
	}
	return a
}
