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
				&dao.Exercise{},
				&dao.Quiz{},
			); err != nil {
				fmt.Printf("migrate database failed:\n%v\n", err)
			}

			fmt.Println("migration success")
		},
	}

	/****************** root cmd *******************/
	var rootCmd = &cobra.Command{Use: "cli"}
	rootCmd.AddCommand(cmdMigrate)
	rootCmd.Execute()
}
