package function

import (
	"encoding/json"
)

// Stringify is a helper function to convert an interface to a string.
func Stringify(i interface{}) string {
	b, _ := json.Marshal(i)
	return string(b)
}

// Parse is a helper function that unfolds a struct into another struct.
// It is important that dest is a pointer to a struct.
func Parse(src interface{}, dest interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dest)
}

// Load is a helper function that unfolds a json string into a struct.
// It is important that dest is a pointer to a struct.
func Load(src string, dest interface{}) error {
	return json.Unmarshal([]byte(src), dest)
}

// LayerMap is a helper function that unfolds a map into another map.
func LayerMap(src map[string]interface{}, dest map[string]interface{}) {
	for k, v := range src {
		dest[k] = v
	}
}
