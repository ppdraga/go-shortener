package withdb

import "time"

type LinkModel struct {
	ID           int64     `gorm:"column:id"`
	Resource     string    `gorm:"column:resource"`
	ShortLink    string    `gorm:"column:short_link"`
	ShortLinkNum int64     `gorm:"column:short_link_num"`
	CustomName   string    `gorm:"column:custom_name"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (LinkModel) TableName() string {
	return "links"
}
