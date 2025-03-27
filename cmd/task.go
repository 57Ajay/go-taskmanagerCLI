package cmd

import (
	"database/sql"
	"fmt"
	"github.com/57ajay/goTask/db"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage your tasks",
	Long:  "Add, list, complete, and remove tasks from your task list.",
}

var taskAddCmd = &cobra.Command{
	Use:   "add [task description]",
	Short: "Add a new task to your list",
	Long:  `Add a new task. The description should be provided as arguments.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Joining all the argumants to form the task description
		description := strings.Join(args, " ")

		sqlStmt := `INSERT INTO tasks (description, status) VALUES (?, ?)`

		// default status to 'pending'
		_, err := db.DB.Exec(sqlStmt, description, "pending")
		if err != nil {
			log.Fatalf("Failed to add task '%s': %v", description, err)
		}

		fmt.Printf("Added task: \"%s\"\n", description)
	},
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your tasks",
	Long:  `Lists all tasks currently stored in your task manager.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Your Tasks:")
		fmt.Println("\n**********-----------**********")
		fmt.Println()
		// Selecting specific columns is generally better than the SELECT *
		sqlStmt := `SELECT id, description, status, created_at, due_date FROM tasks ORDER BY id ASC`

		rows, err := db.DB.Query(sqlStmt)
		if err != nil {
			log.Fatalf("Failed to list tasks: %v", err)
		}
		defer rows.Close()

		taskCount := 0
		for rows.Next() {
			taskCount++
			var id int
			var description, status, createdAt string
			var dueDate sql.NullString

			// Scaning the row data into variables. The, order must match SELECT statement.
			err := rows.Scan(&id, &description, &status, &createdAt, &dueDate)
			if err != nil {
				log.Printf("Warning: Failed to scan row: %v", err)
				continue
			}

			dueDateStr := "N/A"
			if dueDate.Valid {
				dueDateStr = dueDate.String
			}
			fmt.Printf("%d. [%s] %s (Due: %s)\n", id, status, description, dueDateStr)
		}

		if err = rows.Err(); err != nil {
			log.Fatalf("Error iterating over tasks: %v", err)
		}

		if taskCount == 0 {
			fmt.Println("You have no tasks yet!")
		}
		fmt.Println("\n**********-----------**********")
	},
}

var taskDoneCmd = &cobra.Command{
	Use:   "done [task_id]",
	Short: "Mark a task as completed",
	Long:  `Mark a task as completed by providing its ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		idStr := args[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatalf("Invalid task ID provided: '%s'. Please provide a number.", idStr)
		}

		sqlStmt := `UPDATE tasks SET status = ? WHERE id = ?`
		status := "done"

		result, err := db.DB.Exec(sqlStmt, status, id)
		if err != nil {
			log.Fatalf("Failed to mark task %d as done: %v", id, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("Warning: Could not check rows affected after update: %v", err)
			fmt.Printf("Attempted to mark task %d as done.\n", id)
		} else if rowsAffected == 0 {
			fmt.Printf("Task with ID %d not found.\n", id)
		} else {
			fmt.Printf("Marked task %d as done.\n", id)
		}
	},
}

var taskRemoveCmd = &cobra.Command{
	Use:   "remove [task_id]",
	Short: "Remove a task from your list",
	Long:  `Remove a task permanently by providing its ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		idStr := args[0]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatalf("Invalid task ID provided: '%s'. Please provide a number.", idStr)
		}

		sqlStmt := `DELETE FROM tasks WHERE id = ?`

		result, err := db.DB.Exec(sqlStmt, id)
		if err != nil {
			log.Fatalf("Failed to remove task %d: %v", id, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("Warning: Could not check rows affected after delete: %v", err)
			fmt.Printf("Attempted to remove task %d.\n", id)
		} else if rowsAffected == 0 {
			fmt.Printf("Task with ID %d not found.\n", id)
		} else {
			fmt.Printf("Removed task %d.\n", id)
		}
	},
}

func init() {
	AddCommand(taskCmd)
	taskCmd.AddCommand(taskAddCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskDoneCmd)
	taskCmd.AddCommand(taskRemoveCmd)
}
