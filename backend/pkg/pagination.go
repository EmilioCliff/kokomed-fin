package pkg

import "time"

type PaginationMetadata struct {
	PageSize    uint32     `json:"pageSize"`
	CurrentPage uint32     `json:"currentPage"`
	TotalData   uint32     `json:"totalData"`
	TotalPages  uint32     `json:"totalPages"`
	FromDate    *time.Time `json:"from_date,omitempty"`
	ToDate      *time.Time `json:"to_date,omitempty"`
}

func CalculateOffset(currentPage, pageSize uint32) int32 {
	return int32((currentPage - 1) * pageSize)
}

func CreatePaginationMetadata(totalData, pageSize, currentPage uint32) PaginationMetadata {
	return PaginationMetadata{
		PageSize:    pageSize,
		CurrentPage: currentPage,
		TotalData:   totalData,
		TotalPages:  (totalData + pageSize - 1) / pageSize,
	}
}
