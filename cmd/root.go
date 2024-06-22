package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	title       string
	description string
)

var rootCmd = &cobra.Command{
	Use:   "todo",
	Short: "A simple todo application",
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new Todo",
	Long:  "Adds a new Todo item with a title and description",
	Run: func(cmd *cobra.Command, args []string) {
		addTodo(title, description)
	},
}

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read todos",
	Long:  "Read all Todos",
	Run: func(cmd *cobra.Command, args []string) {
		listTodos()
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a todo",
	Long:  `Delete todo item of the id provided. To know the id run "todoify read"`,
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetInt("id")
		if err != nil {
			log.Fatalf("No id provided")
		}

		deleteTodo(id)

	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a todo item",
	Long: `Update a specific todo item identified by ID with new title, description, and IsDone status. To know the id run 
todoify read
  `,
	Run: func(cmd *cobra.Command, args []string) {
		id, _ := cmd.Flags().GetInt("id")
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		isDone, _ := cmd.Flags().GetBool("done")
		fmt.Println(isDone)
		updateTodo(id, title, description, isDone)
	},
}

func init() {
	addCmd.Flags().StringVarP(&title, "title", "t", "", "Add the title of the new Todo")
	addCmd.Flags().StringVarP(&description, "description", "d", "", "Add the description of the new Todo")

	addCmd.MarkFlagRequired("title")
	addCmd.MarkFlagRequired("description")

	deleteCmd.Flags().IntP("id", "i", 0, "Deletes the todo item of the id provided")
	deleteCmd.MarkFlagRequired("id")

	updateCmd.Flags().IntP("id", "i", 0, "ID of the todo to update (required)")
	updateCmd.Flags().StringP("title", "t", "", "New title for the todo")
	updateCmd.Flags().StringP("description", "d", "", "New description for the todo")
	updateCmd.Flags().BoolP("done", "o", false, "Mark todo as done")
	updateCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(updateCmd)

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
