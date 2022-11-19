package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"quizcat/app"
	"quizcat/dao"
)

func main() {

	/****************** migrate cmd *******************/
	var cmdMigrate = &cobra.Command{
		Use:   "migrate",
		Short: "Gorm auto migration",
		Long:  `Gorm auto migration. See More: https://gorm.io/docs/migration.html#Auto-Migration`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := app.DB().AutoMigrate(
				&dao.User{},
				&dao.Token{},
				&dao.Exercise{},
				&dao.Quiz{},
				&dao.Solution{},
				&dao.WordSet{},
				&dao.Word{},
			); err != nil {
				fmt.Printf("migrate database failed:\n%v\n", err)
			}

			fmt.Println("migration success")
		},
	}

	var cmdQuizToSolution = &cobra.Command{
		Use:   "qtos",
		Short: "migrate solutions",
		Long:  `migrate solutions`,
		Run: func(cmd *cobra.Command, args []string) {
			quizToSolution()
		},
	}

	/****************** root cmd *******************/
	var rootCmd = &cobra.Command{Use: "cli"}
	rootCmd.AddCommand(cmdMigrate, cmdQuizToSolution)
	rootCmd.Execute()
}

func quizToSolution() {
	var quizzes []*dao.Quiz
	if err := app.DB().Find(&quizzes).Error; err != nil {
		fmt.Println("get quizzes failed")
	}

	var solutions []*dao.Solution

	for _, q := range quizzes {
		solutions = append(solutions, &dao.Solution{
			Content: q.Solution,
			QuizID:  q.ID,
			UserID:  1,
		})
	}

	if err := app.DB().Create(&solutions).Error; err != nil {
		fmt.Println("create solutions failed")
	}

	fmt.Println("create solutions success.")
}
