package database

import (
	"fmt"
	"strings"

	"cendit.io/garage/primer/enum"
	"cendit.io/garage/primer/typing"
)

// MapToQuery converts an typing.SQLMap to an SQL query string and a slice of arguments
// suitable for use with bun.NewRaw
// if noargs is set to true, the returned slice of arguments will be empty as the value placeholders will be replaced with the actual values
func MapToQuery(m typing.SQLMap, noargs ...bool) (string, []interface{}) {
	var query string
	var args []interface{}
	for k, v := range m.Map {
		if noargs != nil && noargs[0] {
			query += k + " " + string(m.ComparisonOperator) + " " + v.(string) + " " + string(m.JoinOperator) + " "
			continue
		}

		// if value is of type enum.SQLRaw, we need to handle it differently
		if v, ok := v.(enum.SQLRaw); ok {
			query += fmt.Sprintf(" %s %s ", v.Value, m.JoinOperator)
			continue
		}

		// if value is of type enum.SQLAlmostRaw, we need to handle it differently
		if v, ok := v.(enum.SQLAlmostRaw); ok {
			query += fmt.Sprintf(" %s %s %s %s ", k, v.Operator, v.Value, m.JoinOperator)
			continue
		}

		// if value is of type enum.SQLValueMerge, we need to handle it differently
		if v, ok := v.(enum.SQLValueMerge); ok {
			var q string
			if v.Column != "" && v.Operator == "" {
				q = fmt.Sprintf(" %s", v.Column)
			} else {
				if v.Column != "" {
					c := k
					if len(strings.Split(k, ".")) > 1 {
						// prepend table name to column name on the right side of the operator
						c = strings.Split(k, ".")[0] + "." + strings.Split(k, ".")[1]

					}
					q = fmt.Sprintf(" %s %s %s", c, v.Operator, v.Column)
				} else {
					c := k
					if len(strings.Split(k, ".")) > 1 {
						// prepend table name to column name on the right side of the operator
						c = strings.Split(k, ".")[0] + "." + strings.Split(k, ".")[1]
					}
					q = fmt.Sprintf(" %s %s", c, v.Operator)
					for i, val := range v.Values {
						q += " ? "
						if i != len(v.Values)-1 {
							q += string(v.Operator)
						}

						args = append(args, val)
					}
				}
			}
			if len(strings.Split(k, ".")) > 1 {
				// remove table name from column name on the left side of the operator
				k = strings.Split(k, ".")[1]
			}
			query += k + " " + string(m.ComparisonOperator) + q + " " + string(m.JoinOperator) + " "
			continue
		}
		// if comparison operator is enum.In, we need to handle it differently
		if m.ComparisonOperator == enum.In {

			query += k + " " + string(m.ComparisonOperator) + " ("
			for i, val := range v.([]interface{}) {
				query += "?"
				if i != len(v.([]interface{}))-1 {
					query += ","
				}
				args = append(args, val)
			}
			query += ") " + string(m.JoinOperator) + " "
			continue
		}
		query += k + " " + string(m.ComparisonOperator) + " ? " + string(m.JoinOperator) + " "
		args = append(args, v)
	}
	query = query[:len(query)-len(string(m.JoinOperator))-2]
	return query, args
}

// MapsToWQuery converts an typing.SQLMaps to an SQL read query string and a slice of arguments
// suitable for use with bun.NewRaw
func MapsToWQuery(m typing.SQLMaps) (string, []interface{}) {
	var query string
	var args []interface{}
	for _, v := range m.WMaps {
		if len(v.Map) != 0 {
			q, a := MapToQuery(v)
			query += "(" + q + ") " + string(m.WJoinOperator) + " "
			args = append(args, a...)
		}
	}
	if query != "" {
		query = query[:len(query)-len(string(m.WJoinOperator))-2]
	}
	return query, args
}

// MapsToJQuery converts an typing.SQLMaps to an SQL join query string and a slice of arguments
// suitable for use with bun.NewRaw
func MapsToJQuery(m typing.SQLMaps) (string, []interface{}) {
	var query string
	var args []interface{}
	for _, v := range m.JMaps {
		if len(v.Map) != 0 {
			q, a := MapToQuery(v, true)
			query += "(" + q + ") " + string(m.JJoinOperator) + " "
			args = append(args, a...)
		}
	}
	if query != "" {
		query = query[:len(query)-len(string(m.JJoinOperator))-2]
	}
	return query, args
}

// MapsToSQuery converts an typing.SQLMaps to an SQL update query string and a slice of arguments
// suitable for use with bun.NewRaw
func MapsToSQuery(m typing.SQLMaps) (string, []interface{}) {
	wquery, wargs := MapsToWQuery(m)
	squery, sargs := MapToQuery(m.SMap)
	rquery := MapToRQuery(m.RMap)
	q := string(enum.SET) + " " + squery + " "
	if wquery != "" {
		q = q + string(enum.WHERE) + " " + wquery + " "
	}
	if rquery != "" {
		q = q + string(enum.RETURNING) + " " + rquery
	}
	return q, append(sargs, wargs...)
}

// MapToRQuery converts an typing.SQLMap to an SQL returning query string and a slice of arguments
// suitable for use with bun.NewRaw
func MapToRQuery(m typing.SQLMap) string {
	var query string
	var i int
	for k, v := range m.Map {
		if v != nil && m.ComparisonOperator != "" {
			query += k + " " + string(m.ComparisonOperator) + " " + v.(string)
			// if this is the last element, don't add a comma
			if i != len(m.Map)-1 {
				query += ", "
			}
			continue
		}
		query += k
		// if this is the last element, don't add a comma
		if i != len(m.Map)-1 {
			query += ", "
		}
		i++
	}
	return query
}

// MapToSQuery converts an typing.SQLMap to an SQL sum query string
// suitable for use with bun.NewRaw
func MapToSQuery(m typing.SQLMap) string {
	var query string
	i := 0
	for k := range m.Map {
		query += `SUM(` + k + `)`
		// if this is the last element, don't add a comma
		if i != len(m.Map)-1 {
			query += ", "
		}
		i++
	}
	return query
}

// MapToIQuery converts an typing.SQLMap to an SQL insert query string and a slice of arguments
// suitable for use with bun.NewRaw
func MapsToIQuery(m typing.SQLMaps) (string, []interface{}) {
	var query string
	var subquery string
	var args []interface{}

	var aligment []string

	if len(m.IMaps) > 0 {

		if query == "" {
			var i int
			query += "("
			for k := range m.IMaps[0].Map {
				query += k
				aligment = append(aligment, k)
				// if this is the last element, don't add a comma
				if i != len(m.IMaps[0].Map)-1 {
					query += ", "
				}
				i++
			}

			query += ") VALUES "
		}

		j := 0
		for _, _m := range m.IMaps {
			var i int

			for _, v := range aligment {

				// for k, v := range _m.Map {

				// // if the query does not end with "(?, " then add it
				// if !strings.HasSuffix(query, "(?, ") && i == 0 {
				// 	query += "(?"
				// } else {
				// 	query += "?"
				// }

				if i == 0 {
					subquery += "("
				}

				subquery += "?"

				// if this is the last element, don't add a comma
				if i != len(_m.Map)-1 {
					subquery += ", "
				}

				args = append(args, _m.Map[v])
				i++
				// }
			}

			subquery += ")"

			// if this is the last element, don't add a comma
			if j != len(m.IMaps)-1 {
				subquery += ", "
			}

			j++
		}
	}

	query += subquery

	if m.Conflict != nil {

		if len(m.Conflict) > 0 {

			// if no set map is provided, do nothing
			if len(m.SMap.Map) == 0 {
				query += " ON CONFLICT (" + strings.Join(m.Conflict, ", ") + ") DO NOTHING"
				return query, args
			} else {
				query += " ON CONFLICT (" + strings.Join(m.Conflict, ", ") + ") DO UPDATE SET "
				squery, sargs := MapToQuery(m.SMap)
				rquery := MapToRQuery(m.RMap)
				query += squery
				args = append(args, sargs...)
				if rquery != "" {
					query = query + " " + string(enum.RETURNING) + " " + rquery
				}
			}
		} else {

			if len(m.SMap.Map) == 0 {
				query += " ON CONFLICT DO NOTHING"
				return query, args
			} else {
				query += " ON CONFLICT DO UPDATE SET "
				squery, sargs := MapToQuery(m.SMap)
				query += squery
				args = append(args, sargs...)
			}
		}
	}

	return query, args
}

// MapToOQuery converts an typing.SQLMap to an SQL order by query string
// suitable for use with bun.NewRaw
func MapsToOQuery(m typing.SQLMaps) string {
	var query string
	var i int

	for k, v := range m.OMap.Map {
		if !strings.Contains(query, "ORDER BY") {
			query += " ORDER BY "
		}

		query += k + " " + v.(string)
		// if this is the last element, don't add a comma
		if i != len(m.OMap.Map)-1 {
			query += ", "
		}
		i++
	}

	return query
}
