package function

// PInt is a helper function to convert an int to a pointer.
func PInt(i int) *int {
	return &i
}

// PString is a helper function to convert a string to a pointer.
func PString(s string) *string {
	return &s
}

// PBool is a helper function to convert a bool to a pointer.
func PBool(b bool) *bool {
	return &b
}

// PFloat64 is a helper function to convert a float64 to a pointer.
func PFloat64(f float64) *float64 {
	return &f
}

// PAny is a helper function to convert an interface to a pointer.
func PAny(i interface{}) *interface{} {
	return &i
}
