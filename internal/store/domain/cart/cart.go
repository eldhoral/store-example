package cart

const (
    TableName = "cart"
)

type Cart struct {
    ID        int  `json:"id" db:"id"`
    MemberID  int  `json:"member_id" db:"member_id"`
    ProductID int  `json:"product_id" db:"product_id"`
    Quantity  int  `json:"quantity" db:"quantity"`
    IsActive  bool `json:"is_active" db:"is_active"`
}

func (m *Cart) TableName() string {
    return TableName
}
