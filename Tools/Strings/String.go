package Strings

//Contains returns true if the given string is within the array
func Contains(hay []string, needle string) bool {
	for _, a := range hay {
		if a == needle {
			return true
		}
	}
	return false
}
