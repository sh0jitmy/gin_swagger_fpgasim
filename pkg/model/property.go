package model

// model
type Property struct {
	ID    string  `json:"id" example:"frequency"`
	Value string `json:"value" example:"0x00000000"`
	UpdatedAt string `json:"createdAt"`
}
