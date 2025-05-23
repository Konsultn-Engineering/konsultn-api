package builder

import (
	"fmt"
	"gorm.io/gorm"
	"konsultn-api/internal/shared/crud/builder/utils"
	"konsultn-api/internal/shared/crud/types"
	"log"
	"strings"
)

// withFields handles field selection for both regular SELECT and DISTINCT operations
// It processes different types of field specifications and builds the appropriate SQL clauses
// Parameters:
//   - distinct: Whether to use DISTINCT selection
//   - fields: The fields to select, which can be strings, string slices for aliases, or subqueries
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) withFields(distinct bool, fields ...interface{}) types.QueryBuilder[T] {
	var selectCols []string

	for _, field := range fields {
		switch v := field.(type) {
		case string:
			// Simple column
			selectCols = append(selectCols, utils.Quote(v))

		case []string:
			if len(v) == 2 {
				alias := v[1]
				qb.knownAliases[alias] = true
				selectCols = append(selectCols, fmt.Sprintf("%s AS %s", v[0], utils.Quote(v[1])))
			}

		case []interface{}:
			if len(v) == 2 {
				switch sub := v[0].(type) {
				case *QueryBuilder[T]: // Or a generic interface if needed
					rawSQL := fmt.Sprintf("(%s)", sub.ToRawSQL())
					alias, ok := v[1].(string)
					if ok {
						selectCols = append(selectCols, fmt.Sprintf("%s AS %s", rawSQL, utils.Quote(alias)))
					}
				}
			}

		default:
			// Optionally handle or log unsupported types
		}
	}

	if len(selectCols) > 0 {
		if distinct {
			qb.DB = qb.DB.Distinct(strings.Join(selectCols, ","))
		} else {
			qb.DB = qb.DB.Select(strings.Join(selectCols, ", "))
		}
	}

	return qb
}

// Select specifies which fields to retrieve from the database
// Parameters:
//   - fields: Can be field names as strings, or arrays for aliased fields
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) Select(fields ...interface{}) types.QueryBuilder[T] {
	return qb.withFields(false, fields...)
}

// Distinct selects unique records based on the specified fields
// Parameters:
//   - fields: The fields to apply DISTINCT to
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) Distinct(fields ...interface{}) types.QueryBuilder[T] {
	return qb.withFields(true, fields...)
}

// addCondition is a generic method for adding WHERE / OR conditions
// Parameters:
//   - field: The database field name to apply the condition to
//   - value: The value to compare against
//   - operator: The comparison operator (=, >, <, etc.)
//   - isOr: Whether to use OR instead of AND for this condition
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) addCondition(field string, value interface{}, operator Operator, isOr bool) types.QueryBuilder[T] {
	// Quote field name to prevent SQL injection and handle reserved keywords
	quotedField := utils.Quote(field)
	var clause string
	var args []interface{}

	// Check if value is a RawValue
	if rawSQL, rawArgs, isRaw := utils.IsRawValue(value); isRaw {
		// For RawValue, we use the value directly in the SQL statement
		clause = fmt.Sprintf("%s %s %s", quotedField, string(operator), rawSQL)

		// Apply the condition based on OR/AND logic
		if isOr {
			qb.DB = qb.DB.Or(clause, rawArgs...)
		} else {
			qb.DB = qb.DB.Where(clause, rawArgs...)
		}
		return qb
	}

	// Original handling for other value types
	switch operator {
	case IN, NIN:
		// Handle IN/NOT IN operators
		clause = fmt.Sprintf("%s %s (?)", quotedField, string(operator))
		args = []interface{}{value}

	case NULL, NOTNULL:
		// Handle IS NULL / IS NOT NULL (no value needed)
		clause = fmt.Sprintf("%s IS %s", quotedField, string(operator))
		// Empty args slice - no arguments needed

	case BETWEEN:
		// Handle BETWEEN operator with two values
		clause = fmt.Sprintf("%s %s ? AND ?", quotedField, string(operator))

		// Check if value is a slice with exactly two elements
		if slice, ok := value.([]interface{}); ok && len(slice) == 2 {
			args = slice // Use the two elements directly
		} else {
			log.Printf("Error: BETWEEN operator requires exactly two values for field '%s'", field)
			return qb
		}

	default:
		// Handle standard comparison operators (=, >, <, etc.)
		clause = fmt.Sprintf("%s %s ?", quotedField, string(operator))
		args = []interface{}{value}
	}

	// Handle nested query builders (subqueries)
	for i, arg := range args {
		if subQuery, ok := arg.(types.QueryBuilder[T]); ok {
			args[i] = subQuery.G(true)
		}
	}

	// Apply the condition based on OR/AND logic
	if isOr {
		qb.DB = qb.DB.Or(clause, args...)
	} else {
		qb.DB = qb.DB.Where(clause, args...)
	}

	return qb
}

// addWhere is a helper method that adds WHERE conditions with AND logic
// Parameters:
//   - field: The database field name
//   - value: The value to compare against
//   - op: The operator to use for comparison
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) addWhere(field string, value interface{}, op Operator) types.QueryBuilder[T] {
	return qb.addCondition(field, value, op, false)
}

// Where adds an equality condition to the query (field = value)
// Parameters:
//   - field: The database field name
//   - value: The value to compare against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) Where(field string, value interface{}) types.QueryBuilder[T] {
	return qb.addWhere(field, value, EQ)
}

// WhereNull adds a condition checking if the specified field is NULL
// Parameters:
//   - field: The database field name to check for NULL
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) WhereNull(field string) types.QueryBuilder[T] {
	return qb.addWhere(field, nil, NULL)
}

// WhereNotNull adds a condition checking if the specified field is NOT NULL
// Parameters:
//   - field: The database field name to check for NOT NULL
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) WhereNotNull(field string) types.QueryBuilder[T] {
	return qb.addWhere(field, nil, NOTNULL)
}

// WhereNot adds a not-equal condition to the query (field != value)
// Parameters:
//   - field: The database field name
//   - value: The value to compare against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) WhereNot(field string, value interface{}) types.QueryBuilder[T] {
	return qb.addWhere(field, value, NEQ)
}

// WhereLT adds a less-than condition to the query (field < value)
// Parameters:
//   - field: The database field name
//   - value: The value to compare against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) WhereLT(field string, value interface{}) types.QueryBuilder[T] {
	return qb.addWhere(field, value, LT)
}

// WhereLTE adds a less-than-or-equal condition to the query (field <= value)
// Parameters:
//   - field: The database field name
//   - value: The value to compare against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) WhereLTE(field string, value interface{}) types.QueryBuilder[T] {
	return qb.addWhere(field, value, LTE)
}

// WhereGT adds a greater-than condition to the query (field > value)
// Parameters:
//   - field: The database field name
//   - value: The value to compare against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) WhereGT(field string, value interface{}) types.QueryBuilder[T] {
	return qb.addWhere(field, value, GT)
}

// WhereGTE adds a greater-than-or-equal condition to the query (field >= value)
// Parameters:
//   - field: The database field name
//   - value: The value to compare against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) WhereGTE(field string, value interface{}) types.QueryBuilder[T] {
	return qb.addWhere(field, value, GTE)
}

// WhereIN adds a condition checking if the field's value is in a list of values
// Parameters:
//   - field: The database field name
//   - values: A slice of values to check against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) WhereIN(field string, values interface{}) types.QueryBuilder[T] {
	return qb.addWhere(field, values, IN)
}

// WhereNotIN adds a condition checking if the field's value is not in a list of values
// Parameters:
//   - field: The database field name
//   - values: A slice of values to check against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) WhereNotIN(field string, values interface{}) types.QueryBuilder[T] {
	return qb.addWhere(field, values, NIN)
}

// WhereBetween adds a Where BETWEEN value AND another value
// Parameters:
//   - field: The database field name
//   - min: The minimum value for that field
//   - max: The maximum value for that field
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) WhereBetween(field string, min, max interface{}) types.QueryBuilder[T] {
	values := []interface{}{min, max}
	return qb.addWhere(field, values, BETWEEN)
}

// WhereRaw adds a raw SQL condition to the WHERE clause
// This method safely parameterizes any arguments
func (qb *QueryBuilder[T]) WhereRaw(sql string, args ...interface{}) types.QueryBuilder[T] {
	qb.DB = qb.DB.Where(sql, args...)
	return qb
}

// addOrWhere is a helper method that adds WHERE conditions with OR logic
// Parameters:
//   - field: The database field name
//   - value: The value to compare against
//   - op: The operator to use for comparison
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) addOrWhere(field string, value interface{}, op Operator) types.QueryBuilder[T] {
	return qb.addCondition(field, value, op, true)
}

// OrWhere adds an OR equality condition to the query (OR field = value)
// Parameters:
//   - field: The database field name
//   - value: The value to compare against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) OrWhere(field string, value interface{}) types.QueryBuilder[T] {
	return qb.addOrWhere(field, value, EQ)
}

// OrWhereNull adds an OR condition checking if the field is NULL
// Parameters:
//   - field: The database field name to check for NULL
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) OrWhereNull(field string) types.QueryBuilder[T] {
	return qb.addOrWhere(field, nil, NULL)
}

// OrWhereNotNull adds an OR condition checking if the field is NOT NULL
// Parameters:
//   - field: The database field name to check for NOT NULL
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) OrWhereNotNull(field string) types.QueryBuilder[T] {
	return qb.addOrWhere(field, nil, NOTNULL)
}

// OrWhereNot adds an OR not-equal condition to the query (OR field != value)
// Parameters:
//   - field: The database field name
//   - value: The value to compare against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) OrWhereNot(field string, value interface{}) types.QueryBuilder[T] {
	return qb.addOrWhere(field, value, NEQ)
}

// OrWhereGTE adds an OR greater-than-or-equal condition (OR field >= value)
// Parameters:
//   - field: The database field name
//   - value: The value to compare against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) OrWhereGTE(field string, value interface{}) types.QueryBuilder[T] {
	return qb.addOrWhere(field, value, GTE)
}

// OrWhereLTE adds an OR less-than-or-equal condition (OR field <= value)
// Parameters:
//   - field: The database field name
//   - value: The value to compare against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) OrWhereLTE(field string, value interface{}) types.QueryBuilder[T] {
	return qb.addOrWhere(field, value, LTE)
}

// OrWhereIN adds an OR condition checking if field's value is in a list
// Parameters:
//   - field: The database field name
//   - value: A slice of values to check against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) OrWhereIN(field string, value interface{}) types.QueryBuilder[T] {
	return qb.addOrWhere(field, value, IN)
}

// OrWhereNotIN adds an OR condition checking if field's value is not in a list
// Parameters:
//   - field: The database field name
//   - value: A slice of values to check against
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) OrWhereNotIN(field string, value interface{}) types.QueryBuilder[T] {
	return qb.addOrWhere(field, value, NIN)
}

// OrWhereBetween adds a OR Where BETWEEN value AND another value
// Parameters:
//   - field: The database field name
//   - min: The minimum value for that field
//   - max: The maximum value for that field
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) OrWhereBetween(field string, min, max interface{}) types.QueryBuilder[T] {
	values := []interface{}{min, max}
	return qb.addOrWhere(field, values, BETWEEN)
}

// OrWhereRaw adds a raw SQL condition with OR logic
func (qb *QueryBuilder[T]) OrWhereRaw(sql string, args ...interface{}) types.QueryBuilder[T] {
	qb.DB = qb.DB.Or(sql, args...)
	return qb
}

// OrWhereGroup creates a grouped set of OR conditions
// This allows for complex logic like: WHERE ... OR (condition1 AND condition2)
// Parameters:
//   - callback: A function that receives a query builder to define the conditions within the group
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) OrWhereGroup(callback func(types.QueryBuilder[T])) types.QueryBuilder[T] {
	// Create a sub-query builder with a cloned DB for isolation
	subQB := &QueryBuilder[T]{DB: qb.DB.Session(&gorm.Session{NewDB: true})}
	callback(subQB)

	// Apply the sub-condition as a group
	qb.DB = qb.DB.Or(subQB.DB)
	return qb
}

// WhereGroup creates a grouped set of AND conditions
// This allows for complex logic like: WHERE (condition1 AND condition2)
// Parameters:
//   - callback: A function that receives a query builder to define the conditions within the group
//
// Returns: The query builder for method chaining
func (qb *QueryBuilder[T]) WhereGroup(callback func(types.QueryBuilder[T])) types.QueryBuilder[T] {
	subQB := &QueryBuilder[T]{DB: qb.DB.Session(&gorm.Session{NewDB: true})}
	callback(subQB)

	qb.DB = qb.DB.Where(subQB.DB)
	return qb
}
