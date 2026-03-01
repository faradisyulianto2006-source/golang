package entity

type Person struct {
	FirstName string `json:"firstname" binding:"required"`
	LastName string `json:"lastname" binding:"required"`
	Age int `json:"age" binding:"required"`
	Email string `json:"email" validate:"required, email"`
}

// Project struct represents a project entity with its attributes -- Ini mirip model di Laravel
type Project struct {
	// binding tag untuk validasi input, min 3 karakter dan max 100 karakter
	Id int `json:"id" binding:"required" validate:"required" increment"` 
	Title string `json:"title" binding:"min=3,max=100"` 
	Description string `json:"description" binding:"min=10,max=500"`
	Url string `json:"url" binding:"url"`
	Stack []string `json:"stack" binding:"required"`	
}