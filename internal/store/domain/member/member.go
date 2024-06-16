package member

import "time"

const (
    TableName = "member"
)

type Member struct {
    ID          int       `json:"id" db:"id"`
    ChannelID   string    `json:"channel_id" db:"channel_id"`
    Username    string    `json:"username" db:"username"`
    Credential  string    `json:"credential" db:"credential"`
    Salt        string    `json:"salt" db:"salt"`
    CreatedDate time.Time `json:"created_date" db:"created_date"`
}

func (m *Member) TableName() string {
    return TableName
}
