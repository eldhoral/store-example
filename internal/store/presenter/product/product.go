package product

type (
    ProductRequest struct {
        Category string `json:"category"`
    }

    ProductResponse struct {
        ID       int     `json:"id" gorm:"column:id"`
        Name     string  `json:"name" gorm:"column:name"`
        Category string  `json:"category" gorm:"column:category"`
        Price    float64 `json:"price" gorm:"column:price"`
        Stock    int     `json:"stock" gorm:"column:stock"`
    }
)
