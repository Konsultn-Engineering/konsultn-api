package builder

import (
	"fmt"
	"konsultn-api/internal/shared/crud/builder/utils"
	"konsultn-api/internal/shared/crud/types"
	"strings"
)

type SQL struct{}

// Raw creates a raw SQL literal value
// Use for fixed values or SQL expressions without parameters
// WARNING: Never use with user input!
func (s SQL) Raw(value interface{}) types.RawValue {
	return utils.SafeSQL("?", value)
}

// SQL creates a parameterized SQL fragment
// This is the safe way to include dynamic values in SQL
func (s SQL) SQL(sql string, args ...interface{}) types.RawValue {
	return utils.SafeSQL(sql, args...)
}

// Now returns the database NOW() function
func (s SQL) Now() types.RawValue {
	return types.RawValue{Value: "NOW()"}
}

// Count creates a COUNT expression
func (s SQL) Count(field string) types.RawValue {
	return types.RawValue{Value: "COUNT(" + utils.Quote(field) + ")"}
}

// Sum creates a SUM expression
func (s SQL) Sum(field string) types.RawValue {
	return types.RawValue{Value: "SUM(" + utils.Quote(field) + ")"}
}

// Avg creates an AVG expression
func (s SQL) Avg(field string) types.RawValue {
	return types.RawValue{Value: "AVG(" + utils.Quote(field) + ")"}
}

// Min creates a MIN expression
func (s SQL) Min(field string) types.RawValue {
	return types.RawValue{Value: "MIN(" + utils.Quote(field) + ")"}
}

// Max creates a MAX expression
func (s SQL) Max(field string) types.RawValue {
	return types.RawValue{Value: "MAX(" + utils.Quote(field) + ")"}
}

// Coalesce creates a COALESCE expression
func (s SQL) Coalesce(fields []string, defaultValue ...interface{}) types.RawValue {
	sql := "COALESCE(" + strings.Join(fields, ", ")
	if len(defaultValue) > 0 {
		sql += ", ?"
		return utils.SafeSQL(sql+")", defaultValue[0])
	}
	return types.RawValue{Value: sql + ")"}
}

// Cast creates a CAST expression
func (s SQL) Cast(value interface{}, dataType string) types.RawValue {
	return utils.SafeSQL("CAST(? AS "+dataType+")", value)
}

// Lower creates a LOWER function call
func (s SQL) Lower(field string) types.RawValue {
	return types.RawValue{Value: "LOWER(" + utils.Quote(field) + ")"}
}

// Upper creates an UPPER function call
func (s SQL) Upper(field string) types.RawValue {
	return types.RawValue{Value: "UPPER(" + utils.Quote(field) + ")"}
}

// Concat creates a string concatenation expression
func (s SQL) Concat(parts ...interface{}) types.RawValue {
	var sqlParts []string
	var args []interface{}

	for _, part := range parts {
		if str, ok := part.(string); ok {
			// Check if it looks like a column reference
			if !strings.Contains(str, "'") && !strings.Contains(str, "\"") {
				sqlParts = append(sqlParts, utils.Quote(str))
			} else {
				// It's a string literal
				sqlParts = append(sqlParts, "?")
				args = append(args, str)
			}
		} else {
			// Non-string value
			sqlParts = append(sqlParts, "?")
			args = append(args, part)
		}
	}

	return utils.SafeSQL("CONCAT("+strings.Join(sqlParts, ", ")+")", args...)
}

// DateFormat formats a date field (implementation depends on database)
// This version is for PostgreSQL
func (s SQL) DateFormat(field string, format string) types.RawValue {
	return utils.SafeSQL("TO_CHAR("+utils.Quote(field)+", ?)", format)
}

// Interval creates a date interval expression
// e.g., SQL.Interval("1 day") or SQL.Interval("30 minutes")
func (s SQL) Interval(value string) types.RawValue {
	return utils.SafeSQL("INTERVAL ?", value)
}

// Extract extracts a part from a date/time value
// e.g., SQL.Extract("year", "created_at")
func (s SQL) Extract(part string, field string) types.RawValue {
	return types.RawValue{Value: fmt.Sprintf("EXTRACT(%s FROM %s)", part, utils.Quote(field))}
}
