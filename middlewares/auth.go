package middlewares

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"

	"quizcat/app"
	"quizcat/conf"
	"quizcat/dao"
)

func Auth(c *fiber.Ctx) error {
	authToken := c.Get("Authorization")
	if authToken == "" {
		app.Log().Println("Authorization header is not set.")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var token dao.Token

	// get user from cache
	ctx := context.Background()
	prefix := conf.TokenCachePrefix()
	tokenCache, err := app.Cache().Get(ctx, prefix+authToken).Result()
	if err == nil {
		err := json.Unmarshal([]byte(tokenCache), &token)
		if err == nil {
			c.Locals("user", token.User)
			return c.Next()
		}
	}

	if err := app.DB().Preload("User").Where("value = ?", authToken).First(&token).Error; err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusUnauthorized).SendString("load user failed")
	}

	// set token to cache
	tokenBytes, err := json.Marshal(&token)
	if err == nil {
		if err := app.Cache().Set(ctx, prefix+authToken, tokenBytes, time.Hour).Err(); err != nil {
			app.Log().Println(err)
		}
	}

	c.Locals("user", token.User)
	return c.Next()
}
