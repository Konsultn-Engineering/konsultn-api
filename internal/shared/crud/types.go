package crud

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"konsultn-api/internal/shared/crud/pagination"
	"konsultn-api/internal/shared/crud/repository"
	"konsultn-api/internal/shared/crud/types"
)

// Re-export core types
type (
	// Repository expose internal BaseRepository
	Repository[T any, ID comparable] struct {
		*repository.BaseRepository[T, ID]
	}

	// PaginatedResult expose internal PaginatedResult struct
	// PaginatedResult contains paginated query results
	PaginatedResult[T any] = pagination.PaginatedResult[T]
	QueryParams            = pagination.QueryParams

	// IdType defines valid types for primary keys
	IdType = comparable

	Query[T any]        = types.QueryBuilder[T]
	JoinBuilder         = types.JoinBuilder
	QueryBuilder[T any] = types.QueryBuilder[T]
)

func ConvertPaginated[To any](
	from pagination.PaginatedResult[map[string]interface{}], // Explicitly expects map[string]interface{} as 'From'
) (pagination.PaginatedResult[*To], error) {

	var convertedItems []*To

	for i, itemMap := range from.Result {
		var targetItem To // Create a value of type To
		config := &mapstructure.DecoderConfig{
			Metadata:         nil,
			Result:           &targetItem, // No change here, we still need the address
			WeaklyTypedInput: true,
			TagName:          "mapstructure",
		}
		decoder, err := mapstructure.NewDecoder(config)
		if err != nil {
			return pagination.PaginatedResult[*To]{}, fmt.Errorf("failed to create mapstructure decoder: %w", err)
		}

		err = decoder.Decode(itemMap)
		if err != nil {
			// It's often best to return on the first error when doing bulk conversions,
			// or you could collect all errors if partial success is acceptable.
			return pagination.PaginatedResult[*To]{}, fmt.Errorf("failed to decode item %d into %T: %w", i, targetItem, err)
		}
		convertedItems = append(convertedItems, &targetItem)
	}

	// Construct the new paginated result with the converted items and original metadata
	return pagination.PaginatedResult[*To]{
		Result:     convertedItems,
		TotalCount: from.TotalCount,
		Page:       from.Page,
		Limit:      from.Limit,
		TotalPages: from.TotalPages,
		HasNext:    from.HasNext,
		HasPrev:    from.HasPrev,
	}, nil
}
