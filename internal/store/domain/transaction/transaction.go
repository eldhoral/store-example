package transaction

import "time"

const (
    TableName = "transaction"
)

type Transactions struct {
    ID           int       `json:"id" db:"id"`
    MemberID     int       `json:"member_id" db:"member_id"`
    ProductID    int       `json:"product_id" db:"product_id"`
    TrxCode      string    `json:"trx_code" db:"trx_code"`
    ChannelID    string    `json:"channel_id" db:"channel_id"`
    ChannelRefNo string    `json:"channel_ref_no" db:"channel_ref_no"`
    ChannelTime  string    `json:"channel_time" db:"channel_time"`
    ChannelDate  string    `json:"channel_date" db:"channel_date"`
    Amount       float64   `json:"amount" db:"amount"`
    AmountFee    float64   `json:"amount_fee" db:"amount_fee"`
    Status       string    `json:"status" db:"status"`
    Quantity     int       `json:"quantity" db:"quantity"`
    CreatedDate  time.Time `json:"created_date" db:"created_date"`
    UpdatedDate  time.Time `json:"updated_date" db:"updated_date"`
}

func (m *Transactions) TableName() string {
    return TableName
}
