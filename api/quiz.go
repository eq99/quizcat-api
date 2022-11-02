package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"quizcat/app"
	"quizcat/dao"
)

func GetExerciseList(c *fiber.Ctx) error {
	exercises, err := dao.GetExercises()
	if err != nil {
		app.Log().Println(err)
	}

	return c.JSON(exercises)
}

func GetExerciseByID(c *fiber.Ctx) error {
	exerciseID, err := c.ParamsInt("exerciseID")

	if err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "parse param exerciseID error",
		})
	}

	exercise, err := dao.GetExerciseByID(exerciseID)
	if err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "db error",
		})
	}

	return c.JSON(exercise)
}

func GetQuizzesByExerciseID(c *fiber.Ctx) error {
	exID, err := strconv.Atoi(c.Query("exid"))
	if err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "parse query exid error",
		})
	}

	quizzes, err := dao.GetQuizzesByExerciseID(exID)
	if err != nil {
		app.Log().Panicln(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "db error",
		})
	}

	return c.JSON(quizzes)
}
