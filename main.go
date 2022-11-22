package main

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"quizcat/api"
	"quizcat/app"
	m "quizcat/middlewares"
)

func main() {
	router := app.Fiber()
	router.Use(logger.New())
	router.Use(recover.New())

	router.Use(cors.New(cors.Config{
		AllowOrigins: "http://127.0.0.1:5173, http://127.0.0.1:5500, http://127.0.0.1",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// apis
	apiRouter := router.Group("/api/quizcat")

	apiRouter.Get("/exercises", api.GetExerciseList)
	apiRouter.Get("/exercises/:exerciseID", api.GetExerciseByID)
	apiRouter.Get("/quizzes", api.GetQuizzesByExerciseID)
	apiRouter.Post("/solutions", m.Auth, api.GetOrSaveSolution)
	apiRouter.Get("/solutions", api.GetSolutionsByQuizId)

	apiRouter.Get("/ibooks", api.GetInterviewBooks)
	apiRouter.Get("/iquestions", api.GetIQuestionsByBookId)
	apiRouter.Get("/icomments", api.GetICommentsByQuestionId)
	apiRouter.Get("/my/icomments", api.GetICommentsByUserId)
	apiRouter.Post("/icomments", m.Auth, api.SaveIComment)

	apiRouter.Get("/wordsets", api.GetWordSets)
	apiRouter.Get("/wordsets/:setID", api.GetWordSet)

	apiRouter.Post("/captcha", api.SendCaptchaByEmail)
	apiRouter.Post("/signin", api.AuthWithEmail)
	apiRouter.Get("/signout", api.Signout)

	app.GetApp().Run()
}
