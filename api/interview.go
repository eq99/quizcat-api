package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"quizcat/app"
	"quizcat/dao"
)

func GetInterviewBooks(c *fiber.Ctx) error {
	books, err := dao.GetInterviewBooks()
	if err != nil {
		app.Log().Println(err)
	}

	return c.JSON(books)
}

func GetIQuestionsByBookId(c *fiber.Ctx) error {
	bookID, err := strconv.Atoi(c.Query("bookid"))
	if err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "parse bookid faild",
		})
	}

	questions, err := dao.GetIQuestionsByBookId(bookID)
	if err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "get questions failed",
		})
	}

	return c.JSON(questions)
}

func GetICommentsByQuestionId(c *fiber.Ctx) error {
	questionID, err := strconv.Atoi(c.Query("questionid"))
	if err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "parse quesionid faild",
		})
	}

	comments, err := dao.GetICommentsByQuestionId(questionID)
	if err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "get comments failed",
		})
	}

	return c.JSON(comments)
}

func GetICommentsByUserId(c *fiber.Ctx) error {
	token, err := dao.CheckAuthToken(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"msg": "unauthorized",
		})
	}

	bookID, err := strconv.Atoi(c.Query("bookid"))
	if err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "parse bookid faild",
		})
	}

	comments, err := dao.GetICommentsByUserId(token.UserID, bookID)
	if err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "get comments failed",
		})
	}

	return c.JSON(comments)
}

func SaveIComment(c *fiber.Ctx) error {
	user := c.Locals("user").(*dao.User)

	form := &dao.SaveICommentForm{}
	if err := c.BodyParser(form); err != nil {
		app.Log().Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if _, err := dao.SaveIComment(form, user.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"msg": "save comment failed",
		})
	}

	return c.JSON(fiber.Map{
		"msg": "ok",
	})
}
