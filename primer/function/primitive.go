package function

import "strconv"

// StringToInt converts a string to an int
func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
