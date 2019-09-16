package main

// intInSlice checks if int is in slice of ints
func intInSlice(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
