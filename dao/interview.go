package dao

import (
	"time"

	"gorm.io/gorm"

	"quizcat/app"
)

type InterviewBook struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	Name      string         ` gorm:"type:varchar(255)" json:"name"`
	Cover     string         `gorm:"type:varchar(255)" json:"cover"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type IQuestion struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	Body      string         `gorm:"type:text" json:"body"`
	Solution  string         `gorm:"type:text" json:"solution"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	BookID    int            `json:"bookID"`
}

type IComment struct {
	ID          int            `gorm:"primaryKey" json:"id"`
	VoteNum     int            `gorm:"default:0" json:"voteNum"`
	Content     string         `gorm:"type:text" json:"content"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	IQuestionID int            `json:"iQuestionID"`
	UserID      int            `json:"userID"`
	BookID      int            `json:"bookID"`
}

func GetInterviewBooks() ([]*InterviewBook, error) {
	var books []*InterviewBook

	if err := app.DB().Find(&books).Error; err != nil {
		return nil, err
	}

	return books, nil
}

func GetIQuestionsByBookId(bid int) ([]*IQuestion, error) {
	var questions []*IQuestion

	if err := app.DB().
		Where("book_id = ?", bid).
		Find(&questions).
		Error; err != nil {
		return nil, err
	}

	return questions, nil
}

type CommentWithUser struct {
	ID          int       `json:"id"`
	VoteNum     int       `gorm:"default:0" json:"voteNum"`
	Content     string    `gorm:"type:text" json:"content"`
	IQuestionID int       `json:"iQuestionID"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Username    string    `json:"username"`
	Avatar      string    `json:"avatar"`
	UserID      int       `json:"userID"`
}

func GetICommentsByQuestionId(qID int) ([]*CommentWithUser, error) {
	var comments []*CommentWithUser

	if err := app.DB().
		Table("i_comments").
		Where("i_question_id = ?", qID).
		Joins("Join users ON users.id = i_comments.user_id ").
		Select("i_comments.id, i_comments.vote_num, i_comments.content, i_comments.i_question_id, i_comments.updated_at, users.name as username, users.avatar, users.id as user_id").
		Scan(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func GetICommentsByUserId(userID int, bookID int) ([]*IComment, error) {
	var comments []*IComment

	if err := app.DB().
		Where("user_id = ? AND book_id = ?", userID, bookID).
		Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

type SaveICommentForm struct {
	Content     string `validate:"required,min:1,max:4000" form:"content" json:"content"`
	IQuestionID int    `validate:"required" form:"iQuestionID" json:"iQuestionID"`
	BookID      int    `validate:"required" form:"bookID" json:"bookID"`
}

func SaveIComment(form *SaveICommentForm, userID int) (*IComment, error) {
	comment := &IComment{UserID: userID, IQuestionID: form.IQuestionID, BookID: form.BookID}

	if err := app.DB().
		Where("user_id = ? AND i_question_id = ?", userID, form.IQuestionID).
		Assign(IComment{Content: form.Content}).
		FirstOrCreate(&comment).
		Error; err != nil {
		return nil, err
	}

	return comment, nil
}
