package models

type PaginationResponse struct {
	HasNext     bool
	HasPrevious bool
	Count       int64
}

type PaginationRequest struct {
	Page  int
	Limit int
}
