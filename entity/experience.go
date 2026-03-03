package entity

type Experience struct {
	Id int `json:"id" binding:"required" validate:"required" increment"`
	Title string `json:"title" binding:"min=3,max=100"`
	Description string `json:"description" binding:"min=3,max=100"`
	Skills []string `json:"skills" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate string `json:"end_date" binding:"required"`
}