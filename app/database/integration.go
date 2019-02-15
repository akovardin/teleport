package database

import "github.com/jinzhu/gorm"

type Integration struct {
	ID           int    `gorm:"primary_key" json:"id"`
	Title        string `json:"title"`
	Body         string `json:"body"`
	Status       bool   `json:"status"`
	Send         bool   `json:"send"`
	Token        string `json:"token"`
	Channel      string `json:"channel"`
	Secret       string `json:"secret"`
	ProxyAddress string `json:"proxyAddress"`
	ProxyUser    string `json:"proxyUser"`
	ProxyPass    string `json:"proxyPass"`

	UserID int `json:"userId"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func (Integration) TableName() string {
	return "integration"
}

func (p *Integration) FindOne(db *gorm.DB, condition Condition) (*Integration, error) {
	err := find(db.Model(p), condition).First(p).Error

	return p, err
}

func (p *Integration) FindAll(db *gorm.DB, condition Condition) ([]Integration, error) {
	integrations := []Integration{}
	err := find(db.Model(p), condition).Find(&integrations).Error

	return integrations, err
}

func (p *Integration) Count(db *gorm.DB, condition Condition) (int, error) {
	count := 0
	err := find(db.Model(p), condition).Count(&count).Error

	return count, err
}
