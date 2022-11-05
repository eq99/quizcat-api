package api

import (
	"github.com/gofiber/fiber/v2"

	"quizcat/app"
	"quizcat/dao"
)

func GetWordSets(c *fiber.Ctx) error {
	words, err := dao.GetWordSets()
	if err != nil {
		app.Log().Println(err)
	}

	return c.JSON(words)
}

func GetWordSet(c *fiber.Ctx) error {
	setid, err := c.ParamsInt("setID")
	if err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "parse param setID error",
		})
	}

	words, err := dao.GetWordSetByID(setid)
	if err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "db error",
		})
	}

	return c.JSON(words)
}
