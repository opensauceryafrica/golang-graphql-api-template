package primitive

import "strings"

type Array []interface{}

// Len returns the length of the array
func (a Array) Len() int {
	return len(a)
}

// Includes returns true if the array includes the provided value
func (a Array) Includes(val interface{}) bool {
	for _, v := range a {
		if v == val {
			return true
		}
	}
	return false
}

// IndexOf returns the index of the provided value in the array, or -1 if not found
func (a Array) IndexOf(val interface{}) int {
	for i, v := range a {
		if v == val {
			return i
		}
	}
	return -1
}

// ExistsIn reports whether t contains any of the elements of a that is a valid string.
// None string elements in a are ignored during the check.
func (a Array) ExistsIn(t string) bool {
	for _, v := range a {
		if _, ok := v.(string); !ok {
			continue
		}
		if strings.Contains(t, v.(string)) {
			return true
		}
	}
	return false
}
