package dao

import (
	"time"

	"gorm.io/gorm"

	"quizcat/app"
)

type WordSet struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"type:varchar(255)" json:"title"`
	Cover     string         `gorm:"type:varchar(255)" json:"cover"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Word struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	En        string         `gorm:"type:varchar(255)" json:"en"`
	Cn        string         `gorm:"type:varchar(255)" json:"cn"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	WordSetID int            `json:"wordSetID"`
}

func GetWordSets() ([]*WordSet, error) {
	var wordsets []*WordSet

	if err := app.DB().Find(&wordsets).Error; err != nil {
		return nil, err
	}

	return wordsets, nil
}

func GetWordSetByID(id int) ([]*Word, error) {
	var words []*Word

	if err := app.DB().Where("word_set_id = ?", id).Find(&words).Error; err != nil {
		return nil, err
	}

	return words, nil
}
