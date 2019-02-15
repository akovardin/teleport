package database

import "github.com/jinzhu/gorm"

type Integration struct {
	ID          int    `gorm:"primary_key" json:"id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Status      bool   `json:"status"`
	Send        bool   `json:"send"`
	BotKey      string `json:"botKey"`
	ChannelName string `json:"channelName"`
	VKSecret    string `json:"vkSecret"`

	UserID int `json:"userId"`

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func (Integration) TableName() string {
	return "post"
}

func (p Integration) FindOne(db *gorm.DB, condition Condition) (Integration, error) {
	err := Find(db.Model(p), condition).First(&p).Error

	return p, err
}

func (p Integration) FindAll(db *gorm.DB, condition Condition) ([]Integration, error) {
	integrations := []Integration{}
	err := Find(db.Model(p), condition).Find(&integrations).Error

	return integrations, err
}

func (p Integration) Count(db *gorm.DB, condition Condition) (int, error) {
	count := 0
	err := Find(db.Model(p), condition).Count(&count).Error

	return count, err
}
