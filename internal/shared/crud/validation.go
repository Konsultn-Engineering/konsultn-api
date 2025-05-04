package crud

import (
	"fmt"
	"gorm.io/datatypes"
	"reflect"
)

// Validates field types in the update map.
func (m UpdateMap) valid() error {
	for key, value := range m {
		switch value.(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			string,
			datatypes.JSON, datatypes.Date, datatypes.Time,
			datatypes.JSONSlice[any], datatypes.JSONType[any]:
			continue
		default:
			return fmt.Errorf("invalid type for key '%s': %v (type: %s)", key, value, reflect.TypeOf(value).Name())
		}
	}
	return nil
}
