package builder

import (
	"fmt"
	"konsultn-api/internal/shared/crud/builder/utils"
	"konsultn-api/internal/shared/crud/types"
	"strings"
)

// JoinConditionBuilder helps build SQL join conditions (the ON clause)
// It provides a fluent interface for creating complex join conditions
// with AND/OR logic and proper table name quoting
type JoinConditionBuilder struct {
	conditions []string      // The list of individual conditions to join
	params     []interface{} // Parameters for the query
	currentOp  string        // The current operator for joining conditions (AND/OR)
	Base       string        // The base table name or alias
	JoinTable  string        // The joined table name or alias
	Group      bool          // Whether to group the conditions in parentheses
}

// NewJoinConditionBuilder creates a new join condition builder
// Parameters:
//   - base: The base table name or alias (left side of the join)
//   - joinTable: The joined table name or alias (right side of the join)
//
// Returns: A new JoinConditionBuilder instance ready to build conditions
func NewJoinConditionBuilder(base, joinTable string) *JoinConditionBuilder {
	return &JoinConditionBuilder{
		conditions: []string{},
		params:     []interface{}{},
		Base:       base,
		JoinTable:  joinTable,
	}
}

func (j *JoinConditionBuilder) getFieldName(name string, base bool) string {
	if strings.Contains(name, ".") {
		return name
	}

	if base {
		name = utils.Quote(j.Base) + "." + utils.Quote(name)
	} else {
		name = utils.Quote(j.JoinTable) + "." + utils.Quote(name)
	}

	return name
}

// createJoinCondition is an internal helper that formats a single join condition
// It properly quotes table and column names and formats the condition
// Parameters:
//   - left: The column name from the base table
//   - op: The operator for comparison (=, <, >, etc.)
//   - right: The column name from the joined table
func (j *JoinConditionBuilder) createJoinCondition(left, op, right interface{}) {
	var leftStr string
	var rightStr string
	var params []interface{}

	// Process left side
	switch v := left.(type) {
	case string:
		leftStr = j.getFieldName(v, true)
	case types.RawValue:
		leftStr = v.Value
		if v.Args != nil && len(v.Args) > 0 {
			params = append(params, v.Args...)
		}
	default:
		// For other types, use parameterization
		leftStr = "?"
		params = append(params, left)
	}

	// Process right side
	switch v := right.(type) {
	case string:
		rightStr = j.getFieldName(v, false)
	case types.RawValue:
		rightStr = v.Value
		if v.Args != nil && len(v.Args) > 0 {
			params = append(params, v.Args...)
		}
	default:
		// For other types, use parameterization
		rightStr = "?"
		params = append(params, right)
	}

	// Build the condition
	condition := fmt.Sprintf("%s %s %s", leftStr, op, rightStr)
	j.conditions = append(j.conditions, condition)

	// Add parameters
	if len(params) > 0 {
		j.params = append(j.params, params...)
	}
}

// Raw creates a raw SQL fragment that will be included verbatim
// CAUTION: Only use this for fixed SQL fragments, never for user input
func (j *JoinConditionBuilder) Raw(value interface{}) types.RawValue {
	return utils.SafeSQL("?", value)
}

func (j *JoinConditionBuilder) RawSQL(sql string, values ...interface{}) types.RawValue {
	return utils.SafeSQL(sql, values)
}

// On adds an ON condition to the join
// This creates a condition where columns from both tables are compared
// Parameters:
//   - left: The column name from the base table
//   - op: The operator for comparison (typically "=")
//   - right: The column name from the joined table
//
// Example: And("id", "=", "user_id") creates baseTable.id = joinTable.user_id
func (j *JoinConditionBuilder) On(left, op, right interface{}) {
	j.createJoinCondition(left, op, right)
	j.currentOp = "AND"
}

// And adds an AND condition to the join
// This creates a condition where columns from both tables are compared
// Parameters:
//   - left: The column name from the base table
//   - op: The operator for comparison (typically "=")
//   - right: The column name from the joined table
//
// Example: And("id", "=", "user_id") creates baseTable.id = joinTable.user_id
func (j *JoinConditionBuilder) And(left, op, right interface{}) {
	j.On(left, op, right)
}

// Or adds an OR condition to the join
// This creates a condition where either this condition or previous conditions can match
// Parameters:
//   - left: The column name from the base table
//   - op: The operator for comparison (typically "=")
//   - right: The column name from the joined table
//
// Example: Or("secondary_id", "=", "alt_id") creates baseTable.secondary_id = joinTable.alt_id
func (j *JoinConditionBuilder) Or(left, op, right interface{}) {
	j.createJoinCondition(left, op, right)
	j.currentOp = "OR"
}

// String renders the full ON clause condition string
// This converts the built conditions into a SQL string that can be used in a JOIN clause
// If Group is true, the conditions will be wrapped in parentheses
// Returns: A SQL string representing the join conditions
func (j *JoinConditionBuilder) String() string {
	conditions := strings.Join(j.conditions, fmt.Sprintf(" %s ", j.currentOp))
	if j.Group {
		conditions = fmt.Sprintf("(%s)", conditions)
	}
	return conditions
}

// GetParams returns the parameters that should be passed to the query
func (j *JoinConditionBuilder) GetParams() []interface{} {
	return j.params
}
