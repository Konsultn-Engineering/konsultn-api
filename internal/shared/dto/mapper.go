package dto

import (
	"konsultn-api/internal/shared/crud"
)

func MapPaginatedResult[T any, U any](
	source *crud.PaginatedResult[T],
	mapFn func(T) U,
) (PaginatedResultDTO[U], error) {
	// Implementation stays the same
	result := make([]U, len(source.Result))
	for i, item := range source.Result {
		result[i] = mapFn(item)
	}

	return PaginatedResultDTO[U]{
		Result:     result,
		TotalCount: source.TotalCount,
		Page:       source.Page,
		Limit:      source.Limit,
		TotalPages: source.TotalPages,
		HasNext:    source.HasNext,
		HasPrev:    source.HasPrev,
	}, nil
}
