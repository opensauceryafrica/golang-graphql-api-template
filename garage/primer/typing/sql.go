package typing

import "cendit.io/garage/primer/enum"

type SQLMap struct {
	Map map[string]interface{}
	// JoinOperator is the operator that will be used to join each key-value pair in the Map (eg. (name = 'John' AND age = 20))
	JoinOperator enum.SQLOperator
	// ComparisonOperator is the operator that will be used to compare each key-value pair in the Map (eg. name = 'John')
	ComparisonOperator enum.SQLOperator
}

type SQLMaps struct {
	// IMaps for INSERT clauses
	IMaps []SQLMap
	// Conflict for ON CONFLICT clauses (defines the columns that will be used to check for conflicts)
	// use together with IMaps and SMap to define what to do when a conflict is found
	Conflict []string
	// WMaps for WHERE clauses
	WMaps []SQLMap
	// SMaps for SET clauses
	SMap SQLMap
	// RMMap for RETURNING clause
	RMap SQLMap
	// JMaps for JOIN clauses
	JMaps []SQLMap
	// OMap for ORDER BY clauses
	OMap SQLMap
	// Args for the SQL query (these are pre-pended before any other arg already present in the query)
	Args []interface{}
	// WJoinOperator for the SQLMaps present in the WMaps slice. This is used to join the SQLMaps in the WMaps slice (eg. (name = 'John' AND age = 20) OR (name = 'Jane' AND age = 30))
	WJoinOperator enum.SQLOperator
	// JJoinOperator for the SQLMaps present in the JMaps slice. This is used to join the SQLMaps in the JMaps slice (eg. (name = 'John' AND age = 20) OR (name = 'Jane' AND age = 30))
	JJoinOperator enum.SQLOperator
}
