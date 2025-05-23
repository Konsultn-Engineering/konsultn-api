package utils

import (
	"fmt"
	"konsultn-api/internal/shared/crud/types"
	"strings"
)

func Quote(identifier string) string {
	// Don't quote if it looks like a function call (contains parentheses)
	if strings.Contains(identifier, "(") && strings.Contains(identifier, ")") {
		return identifier
	}

	// Basic quoting for PostgreSQL
	parts := strings.Split(identifier, ".")
	for i, p := range parts {
		if p != "*" {
			parts[i] = fmt.Sprintf(`"%s"`, p)
		}
	}
	return strings.Join(parts, ".")
}

func IsRawValue(value interface{}) (string, []interface{}, bool) {
	if rawValue, ok := value.(types.RawValue); ok {
		return rawValue.Value, rawValue.Args, true
	}
	return "", nil, false
}

// SafeSQL ensures SQL fragments are properly parameterized
// It supports both direct strings and parameterized queries
func SafeSQL(sql string, values ...interface{}) types.RawValue {
	placeholderCount := strings.Count(sql, "?")

	// Use exactly the number of values needed for placeholders
	var args []interface{}
	if placeholderCount > 0 && len(values) > 0 {
		// If we have more values than placeholders, only use what we need
		if len(values) > placeholderCount {
			args = values[:placeholderCount]
		} else {
			// Otherwise use all provided values
			args = values
		}
	}

	return types.RawValue{
		Value: sql,
		Args:  args,
	}
}
