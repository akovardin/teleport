package database

import "github.com/jinzhu/gorm"

type Post struct {
	ID        int    `gorm:"primary_key" json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Status    bool   `json:"status"`
	Send      bool   `json:"send"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (Post) TableName() string {
	return "post"
}

func (p Post) FindOne(db *gorm.DB, condition Condition) (Post, error) {
	err := db.Find(db.Model(p), condition).First(&p).Error

	return p, err
}

func (p Post) FindAll(db *gorm.DB, condition Condition) ([]Post, error) {
	posts := []Post{}
	err := Find(db.Model(p), condition).Find(&posts).Error

	return posts, err
}

func (p Post) Count(db *gorm.DB, condition Condition) (int, error) {
	count := 0
	err := Find(db.Model(p), condition).Count(&count).Error

	return count, err
}
