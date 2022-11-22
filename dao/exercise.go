package dao

import (
	"time"

	"gorm.io/gorm"

	"quizcat/app"
)

type Exercise struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	Title     string         ` gorm:"type:varchar(255)" json:"title"`
	QuizNum   int            `gorm:"default:0" json:"quizNum"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Quiz struct {
	ID         int            `gorm:"primaryKey" json:"id"`
	Kind       int            `gorm:"default:0" json:"kind"`
	Level      int            `gorm:"default:1" json:"level"`
	Content    string         `gorm:"type:text" json:"content"`
	Solution   string         `gorm:"type:text" json:"solution"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	ExerciseID int            `json:"exerciseID"`
}

type Solution struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	Score     int            `gorm:"default:0" json:"score"`
	Content   string         `gorm:"type:text" json:"content"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	QuizID    int            `json:"quizID"`
	UserID    int            `json:"userID"`
}

func GetExercises() ([]*Exercise, error) {
	var exercises []*Exercise

	if err := app.DB().Order("created_at desc").Find(&exercises).Error; err != nil {
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

type SaveSolutionForm struct {
	Content string `validate:"required,min:1,max:4000" form:"content" json:"content"`
	QuizID  int    `validate:"required" form:"quizID" json:"quizID"`
}

func GetOrSaveSolution(form *SaveSolutionForm, userID int) (*Solution, error) {
	solution := &Solution{UserID: userID, QuizID: form.QuizID}

	if err := app.DB().
		Where("user_id = ? AND quiz_id = ?", userID, form.QuizID).
		Assign(Solution{Content: form.Content}).
		FirstOrCreate(&solution).
		Error; err != nil {
		return nil, err
	}

	return solution, nil
}

type SolutionWithUser struct {
	ID        int       `json:"id"`
	Score     int       `json:"score"`
	Content   string    `json:"content"`
	QuizID    int       `json:"quizID"`
	UpdatedAt time.Time `json:"updatedAt"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	UserID    int       `json:"userid"`
}

func GetSolutionsByQuizID(quizID int) ([]*SolutionWithUser, error) {
	var solutions []*SolutionWithUser

	if err := app.DB().
		Table("solutions").
		Where("quiz_id = ?", quizID).
		Joins("Join users ON users.id = solutions.user_id ").
		Select("solutions.id, solutions.score, solutions.content, solutions.quiz_id, solutions.updated_at, users.name as username, users.avatar, users.id as user_id").
		Scan(&solutions).Error; err != nil {
		return nil, err
	}

	return solutions, nil
}
