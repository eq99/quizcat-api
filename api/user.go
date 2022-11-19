package api

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"quizcat/app"
	"quizcat/conf"
	"quizcat/dao"
	"quizcat/utils"
)

var ctx = context.Background()

func SendCaptchaByEmail(c *fiber.Ctx) error {
	emailForm := &dao.EmailForm{}

	if err := c.BodyParser(emailForm); err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "parse form failed",
		})
	}

	email := emailForm.Email
	if utils.IsEmailInvalid(email) {
		app.Log().Println("email domain is invalid")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "email domain is invalid",
		})
	}

	captcha := utils.GenCaptcha()
	prefix := conf.CaptchaCachePrefix()
	if err := app.Cache().Set(ctx, prefix+email, captcha, 60*time.Minute).Err(); err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "send email failed",
		})
	}

	if err := utils.SendCaptchaByEmail(captcha, email); err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "send email failed",
		})
	}

	return c.JSON(fiber.Map{
		"msg": "send email success",
	})
}

type AuthWithEmailForm struct {
	Email   string `validate:"required" form:"email" json:"email"`
	Captcha string `validate:"required,max:6" form:"captcha" json:"captcha"`
}

func AuthWithEmail(c *fiber.Ctx) error {
	form := &AuthWithEmailForm{}
	if err := c.BodyParser(form); err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "parse form failed",
		})
	}

	prefix := conf.CaptchaCachePrefix()
	captchaCache, err := app.Cache().Get(ctx, prefix+form.Email).Result()
	if err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "server error",
		})
	}

	if form.Captcha != captchaCache {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "captcha is invalid",
		})
	}

	token, err := dao.AuthUserByEmail(form.Email)
	if err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "server error",
		})
	}

	return c.JSON(token)
}

func Signout(c *fiber.Ctx) error {
	token := c.Get("Authorization")

	if token == "" {
		app.Log().Println("no token in header.")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "not unauthorized",
		})
	}

	if err := dao.Signout(token); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "server error",
		})
	}

	return c.JSON(fiber.Map{
		"msg": "signout success",
	})
}

func PostCheckName(c *fiber.Ctx) error {
	form := &dao.UpdateNameForm{}
	if err := c.BodyParser(form); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if dao.ExistsNickname(form.Name) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"msg": "name is exists",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"msg": "name is ok",
	})
}

func PatchName(c *fiber.Ctx) error {
	user := c.Locals("user").(*dao.User)

	form := &dao.UpdateNameForm{}
	if err := c.BodyParser(form); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if err := dao.UpdateUserName(user.ID, form.Name); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}

type UpdateBioForm struct {
	Bio string `validate:"required,max:40" form:"bio" json:"bio"`
}

func PatchBio(c *fiber.Ctx) error {
	user := c.Locals("user").(*dao.User)

	form := &UpdateBioForm{}
	if err := c.BodyParser(form); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if err := dao.UpdateUserBio(user.ID, form.Bio); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
