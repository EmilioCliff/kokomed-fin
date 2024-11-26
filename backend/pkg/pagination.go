package pkg

type PaginationMetadata struct {
	CurrentPage uint32 `json:"current_page"`
	TotalData   uint32 `json:"total_data"`
	TotalPages  uint32 `json:"total_pages"`
}

const pageSize = 10

func GetPageSize() int32 {
	return pageSize
}

func CalculateOffset(currentPage uint32) int32 {
	return int32((currentPage - 1) * pageSize)
}

func CreatePaginationMetadata(totalData uint32, pageSize uint32, currentPage uint32) PaginationMetadata {
	return PaginationMetadata{
		CurrentPage: currentPage,
		TotalData:   totalData,
		TotalPages:  (totalData + pageSize - 1) / pageSize,
	}
}
