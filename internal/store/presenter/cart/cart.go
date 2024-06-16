package cart

type (
    CartProductDeleteRequest struct {
        MemberID  int `json:"member_id"`
        ProductID int `json:"product_id"`
    }

    CartViewRequest struct {
        MemberID int `json:"member_id"`
    }

    CartRequest struct {
        MemberID  int `json:"member_id" gorm:"column:member_id"`
        ProductID int `json:"product_id" gorm:"column:product_id"`
        Quantity  int `json:"quantity" gorm:"column:quantity"`
    }

    CartResponse struct {
        ID        int  `json:"id" gorm:"column:id"`
        MemberID  int  `json:"member_id" gorm:"column:member_id"`
        ProductID int  `json:"product_id" gorm:"column:product_id"`
        Quantity  int  `json:"quantity" gorm:"column:quantity"`
        IsActive  bool `json:"is_active" gorm:"column:is_active"`
    }
)
