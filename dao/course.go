package dao

import (
	"time"

	"gorm.io/gorm"

	"quizcat/app"
)

type Exercise struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	Title     string         ` gorm:"type:varchar(255)" json:"title"`
	QuizNum   int            `gorm:"default:0" json:"quiz_num"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Quiz struct {
	ID         int            `gorm:"primaryKey" json:"id"`
	Kind       int            `gorm:"default:0" json:"kind"`
	Level      int            `gorm:"default:1" json:"level"`
	Content    string         `gorm:"type:text" json:"content"`
	Solution   string         `gorm:"type:text" json:"solution"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	ExerciseID int            `json:"exercise_id"`
}

func GetExercises() ([]*Exercise, error) {
	var exercises []*Exercise

	if err := app.DB().Find(&exercises).Error; err != nil {
		return nil, err
	}

	return exercises, nil
}

func GetExerciseByID(id int) (*Exercise, error) {
	exercise := &Exercise{ID: id}

	if err := app.DB().First(exercise).Error; err != nil {
		return nil, err
	}

	return exercise, nil
}

func GetQuizzesByExerciseID(eid int) ([]*Quiz, error) {
	var quizzes []*Quiz

	if err := app.DB().Where("exercise_id", eid).Find(&quizzes).Error; err != nil {
		return nil, err
	}

	return quizzes, nil
}
