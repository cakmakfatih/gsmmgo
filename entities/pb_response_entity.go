package entities

type PbListResponseEntity struct {
	Page       int                      `json:"page"`
	PerPage    int                      `json:"perPage"`
	TotalPages int                      `json:"totalPages"`
	TotalItems int                      `json:"totalItems"`
	Items      []map[string]interface{} `json:"items"`
}
