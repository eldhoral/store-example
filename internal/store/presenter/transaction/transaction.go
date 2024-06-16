package transaction

import "time"

type (
    TransactionRequest struct {
        MemberID     int       `json:"member_id" gorm:"column:member_id"`
        ProductID    int       `json:"product_id" gorm:"column:product_id"`
        TrxCode      string    `json:"trx_code" gorm:"column:trx_code"`
        ChannelID    string    `json:"channel_id" gorm:"column:channel_id"`
        ChannelRefNo string    `json:"channel_ref_no" gorm:"column:channel_ref_no"`
        ChannelTime  string    `json:"channel_time" gorm:"column:channel_time"`
        ChannelDate  string    `json:"channel_date" gorm:"column:channel_date"`
        Quantity     int       `json:"quantity" gorm:"column:quantity"`
        CreatedDate  time.Time `json:"created_date" gorm:"column:created_date"`
        UpdatedDate  time.Time `json:"updated_date" gorm:"column:updated_date"`
    }

    TransactionResponse struct {
        ID       int     `json:"id" gorm:"column:id"`
        Name     string  `json:"name" gorm:"column:name"`
        Category string  `json:"category" gorm:"column:category"`
        Price    float64 `json:"price" gorm:"column:price"`
        Stock    int     `json:"stock" gorm:"column:stock"`
    }
)
