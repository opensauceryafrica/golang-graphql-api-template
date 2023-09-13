package function

import (
	"blacheapi/primer/primitive"
	"fmt"
	"reflect"
	"strings"
)

/*
StructTag parses a struct tag and returns the value of the tag

It does this by splitting the tag on commonly used delimiters and returning the first element
in the resulting slice.
*/
func StructTag(tag string) string {
	// common struct tag delimiters
	delimiters := []string{":", ",", " ", "`", "'", "-"}
	for _, v := range delimiters {
		if tag != "" {
			tags := strings.Split(tag, v)
			if len(tags) > 1 {
				return StructTag(tags[0])
			}
		}
	}
	return tag
}

/*
StructToMapOfNonNils uses reflection to convert a struct pointer to a map[string]interface{}
with only non-nil values.

It takes a string to specify which struct tag to use.

It takes an optional slice of slices of strings to specify which fields to include in the map.
If no fields are specified, all fields are included.

It also takes an optional map of string of string to specify which fields to rename in the returned map.
*/
func StructToMapOfNonNils(s interface{}, tag string, fields primitive.Array, replacements map[string]string) map[string]interface{} {
	m := make(map[string]interface{})
	v := reflect.ValueOf(s).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			if fields.Len() > 0 {
				if fields.Includes(StructTag(t.Field(i).Tag.Get(tag))) {
					if replacements != nil && replacements[StructTag(t.Field(i).Tag.Get(tag))] != "" {
						if !v.Field(i).IsZero() && (v.Field(i).Kind() != reflect.Struct) {
							m[replacements[StructTag(t.Field(i).Tag.Get(tag))]] = v.Field(i).Elem().Interface()
							continue
						}
						m[replacements[StructTag(t.Field(i).Tag.Get(tag))]] = v.Field(i).Interface()
						continue
					}
					if !v.Field(i).IsZero() && (v.Field(i).Kind() != reflect.Struct) {
						m[StructTag(t.Field(i).Tag.Get(tag))] = v.Field(i).Elem().Interface()
						continue
					}
					m[StructTag(t.Field(i).Tag.Get(tag))] = v.Field(i).Interface()
				}
				continue
			}
			if replacements != nil && replacements[StructTag(t.Field(i).Tag.Get(tag))] != "" {
				if !v.Field(i).IsZero() && (v.Field(i).Kind() != reflect.Struct) {
					m[replacements[StructTag(t.Field(i).Tag.Get(tag))]] = v.Field(i).Elem().Interface()
					continue
				}
				m[replacements[StructTag(t.Field(i).Tag.Get(tag))]] = v.Field(i).Interface()
				continue
			}
			if !v.Field(i).IsZero() && (v.Field(i).Kind() != reflect.Struct) {
				m[StructTag(t.Field(i).Tag.Get(tag))] = v.Field(i).Elem().Interface()
				continue
			}
			m[StructTag(t.Field(i).Tag.Get(tag))] = v.Field(i).Interface()
		}
	}
	return m
}

/*
ReturnStructFields recursively returns
a slice of interface pointers to the fields of a struct
as need for scanning a database row into a struct
note: the order of the fields in the slice
is the same as the order of the fields in the struct
and this inturn is the same as the order of the columns in the database query

it uses the struct tag `rsf:"false"` to determine fields to skip and the tag `rsfr:"false"` to determine fields to skip on recursive calls
*/
func ReturnStructFields(s interface{}) []interface{} {
	// get the type of argument
	t := reflect.TypeOf(s)
	if t == nil {
		return nil
	}
	// only allow struct type
	if t.Elem().Kind() != reflect.Struct {
		return nil
	}
	// create a slice of interface pointers
	fields := make([]interface{}, t.Elem().NumField())
	// fill the slice with pointers to each struct field
	previousLen := 0
	for i := 0; i < t.Elem().NumField(); i++ {
		// only allow exported fields (`rsf="true"`) or recursive calls
		rsf := t.Elem().Field(i).Tag.Get("rsf") != "false"
		if rsf {
			// if the field is a struct, call this function recursively but ignore the following complex fields
			primitives := primitive.Array{"bun", "time", "mysql"}
			if t.Elem().Field(i).Type.Kind() == reflect.Struct && (t.Elem().Field(i).Tag.Get("rsfr") != "false" && !primitives.ExistsIn(fmt.Sprintf("%T", reflect.New(t.Elem().Field(i).Type).Interface()))) {
				fields = append(fields[:previousLen], append(ReturnStructFields(reflect.ValueOf(s).Elem().Field(i).Addr().Interface()), fields[previousLen:]...)...)
				previousLen += t.Elem().Field(i).Type.NumField()
			} else {
				fields[previousLen] = reflect.ValueOf(s).Elem().Field(i).Addr().Interface()
				previousLen += 1
			}
		}
	}
	// // remove nil values
	for i := 0; i < len(fields); i++ {
		if fields[i] == nil {
			fields = append(fields[:i], fields[i+1:]...)
			i--
		}
	}
	return fields
}

/*
StructFieldToMapOfEqualType uses reflection to select the given fields of a struct pointer and return a map[string]interface{} containing only those fields and their values whose types match the given type

It takes a string to specify which struct tag to use.

It takes an optional slice of slices of strings to specify which fields to include in the map.
If no fields are specified, all fields matching the given type are included.
*/
func StructFieldToMapOfEqualType(s interface{}, tag string, fields primitive.Array, t reflect.Type) map[string]interface{} {
	m := make(map[string]interface{})
	v := reflect.ValueOf(s).Elem()
	tt := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Type() == t {
			if fields.Len() > 0 {
				if fields.Includes(StructTag(tt.Field(i).Tag.Get(tag))) {
					m[StructTag(tt.Field(i).Tag.Get(tag))] = v.Field(i).Interface()
				}
				continue
			}
			m[StructTag(tt.Field(i).Tag.Get(tag))] = v.Field(i).Interface()
		}
	}
	return m
}
