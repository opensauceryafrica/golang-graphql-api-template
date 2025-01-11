package xiao

import (
	"strings"

	"cendit.io/garage/primitive"
)

type SQLMap struct {
	Map map[string]interface{}
	// JoinOperator is the operator that will be used to join each key-value pair in the Map (eg. (name = 'John' AND age = 20))
	JoinOperator SQLOperator
	// ComparisonOperator is the operator that will be used to compare each key-value pair in the Map (eg. name = 'John')
	ComparisonOperator SQLOperator
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
	WJoinOperator SQLOperator
	// JJoinOperator for the SQLMaps present in the JMaps slice. This is used to join the SQLMaps in the JMaps slice (eg. (name = 'John' AND age = 20) OR (name = 'Jane' AND age = 30))
	JJoinOperator SQLOperator
	Pagination    Pagination
}

type Pagination struct {
	// Limit is the number of records to be returned
	Limit int
	// Offset is the number of records to be skipped
	Offset int
}

// SQLOperator is a string type that holds an SQL operator
type SQLOperator string

type SQLDialect string

const (
	// Postgres is the SQL dialect for PostgreSQL
	Postgres SQLDialect = "postgres"
	// MySQL is the SQL dialect for MySQL
	MySQL SQLDialect = "mysql"
)

func (s SQLOperator) String() string {
	return string(s)
}

// SQLValueMerge is a struct that holds an array of values and an operator
// that will be used to merge the values with the value already existent in the column
// Operator & Values work together or you can just use Column to assign the value of a column to another column
type SQLValueMerge struct {
	// the operator to be used to merge the values
	Operator SQLOperator
	// the values to be merged
	Values primitive.Array
	// the column whose value is to be assigned to the operating column
	Column string
}

// SQLAlmostRaw is a struct that holds a raw SQL value and an operator
// the value is assigned as is to the column it is mapped against
type SQLAlmostRaw struct {
	Value    interface{}
	Operator SQLOperator
}

// SQLRaw is a struct that holds a raw SQL value
// the value is used as is and the column it is mapped against is ignored completely in the query
type SQLRaw struct {
	Value interface{}
}

const (
	// Equal is the SQL operator for equality
	Equal SQLOperator = "="

	// NotEqual is the SQL operator for inequality
	NotEqual SQLOperator = "!="

	// GreaterThan is the SQL operator for greater than
	GreaterThan SQLOperator = ">"

	// GreaterThanOrEqual is the SQL operator for greater than or equal to
	GreaterThanOrEqual SQLOperator = ">="

	// LessThan is the SQL operator for less than
	LessThan SQLOperator = "<"

	// LessThanOrEqual is the SQL operator for less than or equal to
	LessThanOrEqual SQLOperator = "<="

	// Like is the SQL operator for the LIKE operator
	Like SQLOperator = "LIKE"

	// NotLike is the SQL operator for the NOT LIKE operator
	NotLike SQLOperator = "NOT LIKE"

	// ILike is the SQL operator for the ILIKE operator
	ILike SQLOperator = "ILIKE"

	// In is the SQL operator for the IN operator
	In SQLOperator = "IN"

	// NotIn is the SQL operator for the NOT IN operator
	NotIn SQLOperator = "NOT IN"

	// IsNull is the SQL operator for the IS NULL operator
	IsNull SQLOperator = "IS NULL"

	// IsNotNull is the SQL operator for the IS NOT NULL operator
	IsNotNull SQLOperator = "IS NOT NULL"

	// Between is the SQL operator for the BETWEEN operator
	Between SQLOperator = "BETWEEN"

	// NotBetween is the SQL operator for the NOT BETWEEN operator
	NotBetween SQLOperator = "NOT BETWEEN"

	// And is the SQL operator for the AND operator
	And SQLOperator = "AND"

	// Or is the SQL operator for the OR operator
	Or SQLOperator = "OR"

	// Not is the SQL operator for the NOT operator
	Not SQLOperator = "NOT"

	// Comma is the SQL operator for the comma
	Comma SQLOperator = ","

	// AS is the SQL operator for AS operator
	AS SQLOperator = "AS"

	// RETURNING is the SQL operator for RETURNING operator
	RETURNING SQLOperator = "RETURNING"

	// SET is the SQL operator for SET operator
	SET SQLOperator = "SET"

	// WHERE is the SQL operator for WHERE operator
	WHERE SQLOperator = "WHERE"

	// PLUS is the SQL operator for PLUS operator
	PLUS SQLOperator = "+"

	// MINUS is the SQL operator for MINUS operator
	MINUS SQLOperator = "-"

	// MULTIPLY is the SQL operator for MULTIPLY operator
	MULTIPLY SQLOperator = "*"

	// DIVIDE is the SQL operator for DIVIDE operator
	DIVIDE SQLOperator = "/"

	// CONCAT is the SQL operator for CONCAT operator
	CONCAT SQLOperator = "||"

	// DESC is the SQL operator for DESC operator
	DESC SQLOperator = "DESC"

	// ASC is the SQL operator for ASC operator
	ASC SQLOperator = "ASC"
)

const (
	// SQL keyword constants
	SQLSelect     SQLOperator = "SELECT"
	SQLInsert     SQLOperator = "INSERT"
	SQLInto       SQLOperator = "INTO"
	SQLUpdate     SQLOperator = "UPDATE"
	SQLDelete     SQLOperator = "DELETE"
	SQLFrom       SQLOperator = "FROM"
	SQLWhere      SQLOperator = "WHERE"
	SQLSet        SQLOperator = "SET"
	SQLValues     SQLOperator = "VALUES"
	SQLReturning  SQLOperator = "RETURNING"
	SQLOnConflict SQLOperator = "ON CONFLICT"
	SQLDoNothing  SQLOperator = "DO NOTHING"
	SQLDoUpdate   SQLOperator = "DO UPDATE"
	SQLAnd        SQLOperator = "AND"
	SQLOr         SQLOperator = "OR"
	SQLIn         SQLOperator = "IN"
	SQLNotIn      SQLOperator = "NOT IN"
	SQLExists     SQLOperator = "EXISTS"
	SQLFor        SQLOperator = "FOR"
	SQLOrderBy    SQLOperator = "ORDER BY"
	SQLAsc        SQLOperator = "ASC"
	SQLDesc       SQLOperator = "DESC"
	SQLBegin      SQLOperator = "BEGIN"
	SQLCommit     SQLOperator = "COMMIT"
	SQLRollback   SQLOperator = "ROLLBACK"
	SQLSum        SQLOperator = "SUM"
)

var (
	SQLCount = func(columns ...string) SQLOperator {
		return SQLOperator("COUNT(" + strings.Join(columns, ", ") + ")")
	}
)
