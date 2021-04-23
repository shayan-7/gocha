package data

type Order struct {
	ID    int    `json:"order_id" validate:"required" gorm:"primaryKey"`
	Title string `json:"title" validate:"required"`
	Price int    `json:"price" validate:"required"`
}
