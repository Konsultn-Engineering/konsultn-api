package crud

import "gorm.io/datatypes"

// idType defines valid types for primary keys.
type idType interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | string
}

// gormDatatype defines GORM-compatible complex field types.
type gormDatatype interface {
	datatypes.JSON | datatypes.Date | datatypes.Time | datatypes.JSONSlice[any] | datatypes.JSONType[any]
}

// UpdateType is a union of acceptable update field types.
type UpdateType interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | string | gormDatatype
}

// UpdateMap defines fields and their new values for updating.
type UpdateMap map[string]interface{}
