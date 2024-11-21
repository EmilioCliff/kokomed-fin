package pkg

type PaginationMetadata struct {
	CurrentPage int32 `json:"current_page"`
	TotalData   int32 `json:"total_data"`
	TotalPages  int32 `json:"total_pages"`
}

const pageSize = 10

func GetPageSize() int32 {
	return pageSize
}

func CalculateOffset(currentPage int32) int32 {
	return (currentPage - 1) * pageSize
}

func CreatePaginationMetadata(totalData int32, pageSize int32, currentPage int32) PaginationMetadata {
	return PaginationMetadata{
		CurrentPage: currentPage,
		TotalData:   totalData,
		TotalPages:  (totalData + pageSize - 1) / pageSize,
	}
}
