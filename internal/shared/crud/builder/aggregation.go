package builder

import (
	"fmt"
	"gorm.io/gorm"
	"konsultn-api/internal/shared/crud/builder/utils"
	"konsultn-api/internal/shared/crud/types"
	"strings"
)

// RawSelect adds a raw SQL expression to the SELECT clause
func (qb *QueryBuilder[T]) RawSelect(expression string, alias string, args ...interface{}) types.QueryBuilder[T] {
	if alias != "" {
		expression = fmt.Sprintf("%s AS %s", expression, utils.Quote(alias))
		qb.knownAliases[alias] = true
	}

	// Use GORM's ability to handle parameters in Select
	qb.DB = qb.DB.Select(gorm.Expr(expression, args...))
	return qb
}

// SQL Function Helpers

// Now returns a RawValue with the database's NOW() function
func (qb *QueryBuilder[T]) Now() types.RawValue {
	return qb.Raw("NOW()")
}

// Cast creates a parameterized CAST expression
func (qb *QueryBuilder[T]) Cast(value interface{}, dataType string) types.RawValue {
	return qb.Raw("CAST(? AS "+dataType+")", value)
}

// Coalesce creates a parameterized COALESCE expression
func (qb *QueryBuilder[T]) Coalesce(fields []string, defaultValue ...interface{}) types.RawValue {
	var args []interface{}
	var placeholders []string

	// Convert field names to proper SQL
	for _, field := range fields {
		if strings.Contains(field, ".") {
			// Field name with table qualifier
			placeholders = append(placeholders, field)
		} else {
			// Simple field name
			placeholders = append(placeholders, utils.Quote(field))
		}
	}

	// Add default value if provided
	if len(defaultValue) > 0 {
		placeholders = append(placeholders, "?")
		args = append(args, defaultValue[0])
	}

	return qb.Raw("COALESCE("+strings.Join(placeholders, ", ")+")", args...)
}

// GroupBy adds a GROUP BY clause to the query
// This function is used to group rows that have the same values in specified columns
// into aggregated results
// Parameters:
//   - fields: One or more field names to group by
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) GroupBy(fields ...string) types.QueryBuilder[T] {
	for i, field := range fields {
		fields[i] = utils.Quote(field)
	}
	qb.DB = qb.DB.Group(strings.Join(fields, ", "))
	return qb
}

// Having adds a HAVING condition to the query (for use with GROUP BY)
// HAVING filters the results of GROUP BY aggregations similar to how WHERE filters individual rows
// Parameters:
//   - condition: Raw SQL condition string
//   - args: Values for any placeholders in the condition
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) Having(condition string, args ...interface{}) types.QueryBuilder[T] {
	qb.DB = qb.DB.Having(condition, args...)
	return qb
}

// HavingField adds a field-based HAVING condition using the provided operator
// This is a helper method used by specific Having methods that handles the SQL formatting
// Parameters:
//   - field: The field name to apply the condition to
//   - value: The value to compare against
//   - op: The operator to use for comparison
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) HavingField(field string, value interface{}, op Operator) types.QueryBuilder[T] {
	clause := fmt.Sprintf("%s %s ?", utils.Quote(field), string(op))
	qb.DB = qb.DB.Having(clause, value)
	return qb
}

// HavingEQ adds a HAVING condition checking for equality (field = value)
// Used to filter groups where an aggregate value equals the specified value
// Parameters:
//   - field: The field name to apply the condition to
//   - value: The value to compare against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) HavingEQ(field string, value interface{}) types.QueryBuilder[T] {
	return qb.HavingField(field, value, EQ)
}

// HavingNEQ adds a HAVING condition checking for inequality (field != value)
// Used to filter groups where an aggregate value does not equal the specified value
// Parameters:
//   - field: The field name to apply the condition to
//   - value: The value to compare against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) HavingNEQ(field string, value interface{}) types.QueryBuilder[T] {
	return qb.HavingField(field, value, NEQ)
}

// HavingGT adds a HAVING condition checking for greater than (field > value)
// Used to filter groups where an aggregate value is greater than the specified value
// Parameters:
//   - field: The field name to apply the condition to
//   - value: The value to compare against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) HavingGT(field string, value interface{}) types.QueryBuilder[T] {
	return qb.HavingField(field, value, GT)
}

// HavingGTE adds a HAVING condition checking for greater than or equal (field >= value)
// Used to filter groups where an aggregate value is greater than or equal to the specified value
// Parameters:
//   - field: The field name to apply the condition to
//   - value: The value to compare against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) HavingGTE(field string, value interface{}) types.QueryBuilder[T] {
	return qb.HavingField(field, value, GTE)
}

// HavingLT adds a HAVING condition checking for less than (field < value)
// Used to filter groups where an aggregate value is less than the specified value
// Parameters:
//   - field: The field name to apply the condition to
//   - value: The value to compare against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) HavingLT(field string, value interface{}) types.QueryBuilder[T] {
	return qb.HavingField(field, value, LT)
}

// HavingLTE adds a HAVING condition checking for less than or equal (field <= value)
// Used to filter groups where an aggregate value is less than or equal to the specified value
// Parameters:
//   - field: The field name to apply the condition to
//   - value: The value to compare against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) HavingLTE(field string, value interface{}) types.QueryBuilder[T] {
	return qb.HavingField(field, value, LTE)
}

// HavingIN adds a HAVING condition checking if a field's value is in a list of values
// Used to filter groups where an aggregate value matches any value in the provided list
// Parameters:
//   - field: The field name to apply the condition to
//   - values: A slice of values to check against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) HavingIN(field string, values interface{}) types.QueryBuilder[T] {
	clause := fmt.Sprintf("%s IN (?)", utils.Quote(field))
	qb.DB = qb.DB.Having(clause, values)
	return qb
}

// HavingBetween adds a HAVING condition checking if a field's value between min and max
// Parameters:
//   - field: The field name to apply the condition to
//   - values: A slice of values to check against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) HavingBetween(field string, min, max interface{}) types.QueryBuilder[T] {
	clause := fmt.Sprintf("%s BETWEEN ? AND ?", utils.Quote(field))
	qb.DB = qb.DB.Having(clause, min, max)
	return qb
}

// OrHaving adds an OR HAVING condition
// This allows multiple HAVING conditions to be combined with OR logic
// Parameters:
//   - condition: Raw SQL condition string
//   - args: Values for any placeholders in the condition
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) OrHaving(condition string, args ...interface{}) types.QueryBuilder[T] {
	qb.DB = qb.DB.Or(condition, args...).Having("")
	return qb
}

// HavingGroup adds a grouped HAVING condition
// This allows for complex HAVING conditions with nested AND/OR logic
// Parameters:
//   - callback: A function that receives a query builder to define the conditions within the group
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) HavingGroup(callback func(types.QueryBuilder[T])) types.QueryBuilder[T] {
	subQB := &QueryBuilder[T]{DB: qb.DB.Session(&gorm.Session{NewDB: true})}
	callback(subQB)

	// Use Having with the subquery
	qb.DB = qb.DB.Having(subQB.DB)
	return qb
}

// OrHavingGroup adds a grouped OR HAVING condition
// This allows for complex HAVING conditions with nested AND/OR logic combined with OR
// Parameters:
//   - callback: A function that receives a query builder to define the conditions within the group
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) OrHavingGroup(callback func(types.QueryBuilder[T])) types.QueryBuilder[T] {
	subQB := &QueryBuilder[T]{DB: qb.DB.Session(&gorm.Session{NewDB: true})}
	callback(subQB)

	// Use Or with the subquery in Having context
	qb.DB = qb.DB.Or(subQB.DB).Having("")
	return qb
}

// Count returns the number of records that match the current query conditions
// This is an aggregation function that returns the total count without retrieving all records
// Returns:
//   - int64: The number of matching records
//   - error: Any error that occurred during counting
func (qb *QueryBuilder[T]) Count() (int64, error) {
	var count int64
	err := qb.build().Count(&count).Error
	return count, err
}

// Exists checks if any records match the current query conditions
// This is a convenience method that uses Count internally but returns a boolean result
// Returns:
//   - bool: True if at least one record matches, false otherwise
//   - error: Any error that occurred during the check
func (qb *QueryBuilder[T]) Exists() (bool, error) {
	count, err := qb.Count()
	return count > 0, err
}
