package models

type Pet struct {
	Name     string `json:"name"`
	Id       int64  `json:"id"`
	Category string `json:"category"`
	Notes    string `json:"notes"`
}

type CategoryResponse struct {
	Pets []Pet `json:"data"`
}
