package dao

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"quizcat/app"
	"quizcat/conf"
)

type User struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(64)" json:"name"`
	Avatar    string         `gorm:"type:varchar(255)" json:"avatar"`
	Email     string         `gorm:"type:varchar(255)" json:"email"`
	IsAdmin   bool           `gorm:"default:false" json:"isAdmin"`
	Status    int            `gorm:"default:0" json:"status"` // 0: normal, 1: blocked
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Token struct {
	ID        int            `gorm:"primaryKey" json:"id"`
	Value     uuid.UUID      `gorm:"index;type:uuid" json:"value"`
	Client    string         `json:"client"`
	UserID    int            `gorm:"unique" json:"userId"` // in this version, only one token is ok
	User      *User          `json:"user"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func DeleteToken(token string) ([]*Token, error) {
	var tokens []*Token
	err := app.DB().Clauses(clause.Returning{}).Where("token = ?", token).Delete(&tokens).Error
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func GetTokenByUserId(userID int) (*Token, error) {
	token := &Token{}
	if err := app.DB().
		Model(token).
		Where("user_id = ?", userID).
		First(&token).Error; err != nil {
		return nil, err
	}

	return token, nil
}

var ctx = context.Background()

func CheckAuthToken(c *fiber.Ctx) (*Token, error) {
	tokenStr := c.Get("Authorization")

	if tokenStr == "" {
		app.Log().Println("no token in header.")
		return nil, fiber.ErrUnauthorized
	}

	// get token from cache
	var token Token
	cachePrefix := conf.TokenCachePrefix()
	tokenCache, err := app.Cache().Get(ctx, cachePrefix+tokenStr).Result()
	if err == nil {
		err := json.Unmarshal([]byte(tokenCache), &token)
		if err == nil {
			return &token, nil
		}
	}

	if err := app.DB().Preload("User").Where("value = ?", tokenStr).First(&token).Error; err != nil {
		app.Log().Println(err)
		return nil, fiber.ErrUnauthorized
	}

	// set user to cache
	tokenBytes, err := json.Marshal(&token)
	if err == nil {
		if err := app.Cache().Set(ctx, cachePrefix+tokenStr, tokenBytes, time.Hour).Err(); err != nil {
			app.Log().Println(err)
		}
	}

	return &token, nil
}

func AuthUserByEmail(email string) (*Token, error) {
	name := uuid.New().String() // use uuid as default name
	user := &User{
		Email:  email,
		Name:   name,
		Avatar: "https://api.multiavatar.com/" + strings.Split(name, "-")[0] + ".svg",
	}

	// get or create user
	result := app.DB().Where(User{Email: email}).FirstOrCreate(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	// update or create token for user
	// this requires user_id is unique.
	token := &Token{UserID: user.ID, Value: uuid.New()}
	if err := app.DB().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(token).Error; err != nil {
		return nil, err
	}

	token.User = user
	return token, nil
}

func CreateAdminUser(email, name, avatar string) (*User, error) {
	user := &User{
		Email:   email,
		Name:    name,
		Avatar:  "https://api.multiavatar.com/" + name + ".svg",
		IsAdmin: true,
	}

	if err := app.DB().Where(User{Email: email}).FirstOrCreate(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func ExistsNickname(name string) bool {
	exists := true
	err := app.DB().Model(&User{}).
		Select("count(*) > 0").
		Where("name = ?", name).
		Find(&exists).
		Error

	if err != nil {
		return true
	}

	return exists
}

func GetUserByEmail(email string) (*User, error) {
	user := &User{Email: email}
	if err := app.DB().Where(&User{Email: email}).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUserName(userID int, newName string) error {
	return app.DB().Model(&User{}).Where("id = ?", userID).Update("name", newName).Error
}

func UpdateUserBio(userID int, newBio string) error {
	return app.DB().Model(&User{}).Where("id = ?", userID).Update("bio", newBio).Error
}
