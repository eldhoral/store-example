package product

const (
    TableName = "product"
)

type Product struct {
    ID       int     `json:"id" db:"id"`
    Name     string  `json:"name" db:"name"`
    Category string  `json:"category" db:"category"`
    Price    float64 `json:"price" db:"price"`
    Stock    int     `json:"stock" db:"stock"`
}

func (m *Product) TableName() string {
    return TableName
}
