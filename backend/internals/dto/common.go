package dto

type BadResponse struct {
	Message string `json:"message"`
}

type Meta struct {
	Page    uint `json:"page"`
	PerPage uint `json:"per_page"`
	Total   uint `json:"total"`
}
