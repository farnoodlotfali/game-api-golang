package models

type Pagination struct {
	Total       int64 `json:"total"`
	LastPage    int   `json:"last_page"`
	CurrentPage int   `json:"current_page"`
}

type PaginationType[T any] struct {
	Items      T          `json:"items"`
	Pagination Pagination `json:"pagination"`
}

type PageResponseType[T any] struct {
	Message string            `json:"message"`
	Data    PaginationType[T] `json:"data"`
}

// normal
type ResponseType[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func SuccessPaginationResponse[T any](items []T, total int64, lastPage, intPage int) (PageResponseType[[]T], error) {
	return PageResponseType[[]T]{
		Message: "nice",
		Data: PaginationType[[]T]{
			Items: items,
			Pagination: Pagination{
				Total:       total,
				LastPage:    lastPage,
				CurrentPage: intPage,
			},
		},
	}, nil
}
