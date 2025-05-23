package builder

import (
	"fmt"
	"gorm.io/gorm"
	"konsultn-api/internal/shared/crud/builder/utils"
	"konsultn-api/internal/shared/crud/types"
	"strings"
)

// JoinClause represents a SQL JOIN operation with its type, table, and conditions
// It contains all the information needed to build a complete JOIN clause
type JoinClause struct {
	Type      JoinType              // The type of join (INNER, LEFT, RIGHT, CROSS)
	Table     string                // The table name or alias to join with
	Condition *JoinConditionBuilder // The join conditions (ON clause)
}

// Join performs an INNER JOIN with the specified table
// This is the standard join that returns rows when there is at least one match in both tables
// Parameters:
//   - table: The table name to join
//   - opts: Optional parameters, where the first element can be the table alias
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) Join(table string, opts ...string) types.QueryBuilder[T] {
	return qb.addJoin(table, opts, JoinInner)
}

// LeftJoin performs a LEFT JOIN with the specified table
// Returns all rows from the left table and matched rows from the right table
// Parameters:
//   - table: The table name to join
//   - opts: Optional parameters, where the first element can be the table alias
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) LeftJoin(table string, opts ...string) types.QueryBuilder[T] {
	return qb.addJoin(table, opts, JoinLeft)
}

// RightJoin performs a RIGHT JOIN with the specified table
// Returns all rows from the right table and matched rows from the left table
// Parameters:
//   - table: The table name to join
//   - opts: Optional parameters, where the first element can be the table alias
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) RightJoin(table string, opts ...string) types.QueryBuilder[T] {
	return qb.addJoin(table, opts, JoinRight)
}

// CrossJoin performs a CROSS JOIN with the specified table
// Returns the Cartesian product of both tables (all possible combinations of rows)
// Parameters:
//   - table: The table name to join
//   - opts: Optional parameters, where the first element can be the table alias
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) CrossJoin(table string, opts ...string) types.QueryBuilder[T] {
	return qb.addJoin(table, opts, JoinCross)
}

// addJoin is a helper method that implements the common functionality for all join types
// It handles table name quoting, aliasing, and adds the join to the internal list
// Parameters:
//   - table: The table name to join
//   - opts: Optional parameters for aliasing the table
//   - joinType: The type of join to perform (INNER, LEFT, RIGHT, CROSS)
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) addJoin(table string, opts []string, joinType JoinType) types.QueryBuilder[T] {
	alias := table
	if len(opts) > 0 {
		alias = opts[0]
		table = fmt.Sprintf("%s AS %s", utils.Quote(table), utils.Quote(alias))
	} else {
		table = utils.Quote(table)
	}

	qb.lastJoinedTable = alias // Track only alias (or original name if no alias)

	qb.joins = append(qb.joins, JoinClause{
		Type:      joinType,
		Table:     table,
		Condition: NewJoinConditionBuilder(qb.baseTable, alias),
	})

	return qb
}

// On specifies the join condition between tables using an equality comparison
// This method should be called immediately after a join method
// Parameters:
//   - left: The column name from the base or previously joined table
//   - right: The column name from the newly joined table
//
// Returns: The query builder for method chaining
// Panics: If called without a preceding join operation
func (qb *QueryBuilder[T]) On(left, right string) types.QueryBuilder[T] {
	if len(qb.joins) == 0 {
		panic("No join to apply condition to")
	}

	lastJoin := &qb.joins[len(qb.joins)-1]
	lastJoin.Condition.And(left, "=", right)
	return qb
}

// OnGroup allows creating complex join conditions with a builder function
// This enables conditions with multiple parts, custom operators, and OR logic
// Parameters:
//   - builder: A function that receives a join condition builder to define complex join conditions
//
// Returns: The query builder for method chaining
// Panics: If called without a preceding join operation
func (qb *QueryBuilder[T]) OnGroup(builder func(j types.JoinBuilder)) types.QueryBuilder[T] {
	if len(qb.joins) == 0 {
		panic("No join to apply condition to")
	}
	lastJoin := &qb.joins[len(qb.joins)-1]
	lastJoin.Condition.Group = true
	builder(lastJoin.Condition)
	return qb
}

// buildJoins constructs the SQL for all joins and applies them to the query
// This is an internal method used when finalizing the query before execution
// Returns: The modified gorm.DB instance with all join clauses applied
func (qb *QueryBuilder[T]) buildJoins() *gorm.DB {
	db := qb.DB

	for _, join := range qb.joins {
		joinSQL := fmt.Sprintf("%s %s", strings.ToUpper(join.Type.String()), join.Table)

		if conditionSQL := join.Condition.String(); conditionSQL != "" {
			joinSQL += " ON " + conditionSQL
		}

		// Get parameters for this join condition
		params := join.Condition.GetParams()

		if len(params) > 0 {
			// If we have parameters, use them with the join
			db = db.Joins(joinSQL, params...)
		} else {
			// Otherwise do a simple join without parameters
			db = db.Joins(joinSQL)
		}
	}

	return db
}
