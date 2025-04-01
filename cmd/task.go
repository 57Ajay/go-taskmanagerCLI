package cmd

import (
	"database/sql"
	"fmt"
	"github.com/57ajay/goTask/db"
	"github.com/spf13/cobra"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var dueDateStr string
var durationRegex = regexp.MustCompile(`(\d+)\s*(y|M|w|d|h|m|s)`)
var silentIfClear bool

func parseRelativeDuration(durationStr string) (time.Time, error) {
	now := time.Now()
	targetTime := now

	years, months, days := 0, 0, 0
	var timeDuration time.Duration

	matches := durationRegex.FindAllStringSubmatch(durationStr, -1)

	if len(matches) > 0 {
		parsedAny := false
		remainingStr := durationStr

		for _, match := range matches {
			if len(match) != 3 {
				continue
			}

			valueStr, unit := match[1], match[2]
			remainingStr = strings.Replace(remainingStr, match[0], "", 1)

			value, err := strconv.Atoi(valueStr)
			if err != nil {
				return time.Time{}, fmt.Errorf("invalid number '%s' found by regex", valueStr)
			}
			parsedAny = true

			switch unit {
			case "y":
				years += value
			case "M": // Month
				months += value
			case "w":
				days += value * 7
			case "d":
				days += value
			case "h":
				timeDuration += time.Duration(value) * time.Hour
			case "m":
				timeDuration += time.Duration(value) * time.Minute
			case "s":
				timeDuration += time.Duration(value) * time.Second
			}
		}

		remainingStr = strings.TrimSpace(remainingStr)
		if remainingStr != "" {
			d, err := time.ParseDuration(remainingStr)
			if err == nil {
				timeDuration += d
				parsedAny = true
			} else {
				return time.Time{}, fmt.Errorf("invalid or unparsed components in duration: '%s'", remainingStr)
			}
		}

		if !parsedAny {
			return time.Time{}, fmt.Errorf("no valid duration components processed")
		}

		targetTime = targetTime.AddDate(years, months, days)
		targetTime = targetTime.Add(timeDuration)

		return targetTime, nil

	} else {
		d, err := time.ParseDuration(durationStr)
		if err == nil {
			return targetTime.Add(d), nil
		}
	}
	return time.Time{}, fmt.Errorf("string is not a recognized relative duration")
}

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Manage your tasks",
	Long:  "Add, list, complete, and remove tasks from your task list.",
}

var taskAddCmd = &cobra.Command{
	Use:   "add [task description]",
	Short: "Add a new task to your list",
	Long:  `Add a new task. The description is required. Use --due flag for due date.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		description := strings.Join(args, " ")
		status := "pending"

		var sqlStmt string
		var err error
		var parsedDueDate any

		if dueDateStr != "" {
			var t time.Time
			var parseErr error

			t, parseErr = parseRelativeDuration(dueDateStr)

			if parseErr != nil {
				layouts := []string{
					"2006-01-02 15:04:05",
					"2006-01-02",
				}
				isAbsolute := false
				for _, layout := range layouts {
					t, parseErr = time.Parse(layout, dueDateStr)
					if parseErr == nil {
						isAbsolute = true
						break
					}
				}
				if !isAbsolute {
					log.Printf("Warning: Invalid due date format: '%s'. Use 'YYYY-MM-DD', 'YYYY-MM-DD HH:MM:SS', or relative like '2d', '3h30m'. Task added without due date.", dueDateStr)
					parsedDueDate = nil
					parseErr = nil
				}
			}

			if parseErr == nil && parsedDueDate == nil {
				if !t.IsZero() {
					parsedDueDate = t.Format("2006-01-02 15:04:05")
					fmt.Printf("Due date set: %s\n", parsedDueDate)
				} else if dueDateStr != "" {
					parsedDueDate = nil
				}
			}

		} else {
			parsedDueDate = nil
		}

		if parsedDueDate != nil {
			sqlStmt = `INSERT INTO tasks (description, status, due_date) VALUES (?, ?, ?)`
			_, err = db.DB.Exec(sqlStmt, description, status, parsedDueDate)
		} else {
			sqlStmt = `INSERT INTO tasks (description, status) VALUES (?, ?)`
			_, err = db.DB.Exec(sqlStmt, description, status)
		}

		if err != nil {
			log.Fatalf("Failed to add task '%s': %v", description, err)
		}

		fmt.Printf("Added task: \"%s\"\n", description)

		dueDateStr = ""
	},
}

// cmd/task.go (inside taskListCmd definition)
var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your tasks",
	Long:  `Lists all tasks currently stored, indicating status and due date. Overdue tasks are marked.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Your Tasks:")
		fmt.Println("-----------")

		sqlStmt := `SELECT id, description, status, created_at, due_date FROM tasks ORDER BY id ASC`
		rows, err := db.DB.Query(sqlStmt)
		if err != nil {
			log.Fatalf("Failed to list tasks: %v", err)
		}
		defer rows.Close()

		taskCount := 0
		now := time.Now() // Get current time once for comparison

		for rows.Next() {
			taskCount++
			var id int
			var description, status, createdAt string
			var dueDateDB sql.NullString

			err := rows.Scan(&id, &description, &status, &createdAt, &dueDateDB)
			if err != nil {
				log.Printf("Warning: Failed to scan row: %v", err)
				continue
			}

			dueDateDisplay := "N/A"
			isOverdue := false
			// var dueDateTime time.Time

			if dueDateDB.Valid {
				parsedTime, parseErr := time.Parse("2006-01-02 15:04:05", dueDateDB.String)
				if parseErr == nil {
					// dueDateTime := parsedTime
					dueDateDisplay = parsedTime.Format("Mon, 02 Jan 2006 15:04")
					if status != "done" && parsedTime.Before(now) {
						isOverdue = true
					}
				} else {
					dueDateDisplay = "Invalid Date Format in DB"
					log.Printf("Warning: Could not parse due date '%s' from DB for task ID %d: %v", dueDateDB.String, id, parseErr)
				}
			}

			overdueMarker := ""
			if isOverdue {
				overdueMarker = " [OVERDUE!]"
			}

			statusMarker := fmt.Sprintf("[%s]", status)
			if status == "done" {
				statusMarker = "[DONE]"
			}

			fmt.Printf("%d. %-9s %s (Due: %s)%s\n", id, statusMarker, description, dueDateDisplay, overdueMarker)
		}

		if err = rows.Err(); err != nil {
			log.Fatalf("Error iterating over tasks: %v", err)
		}

		if taskCount == 0 {
			fmt.Println("You have no tasks yet!")
		}
		fmt.Println("-----------")
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

var taskCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for overdue and upcoming tasks",
	Long:  `Scans your tasks and displays any that are overdue or due within the next 24 hours.`,
	Run: func(cmd *cobra.Command, args []string) {

		now := time.Now()
		upcomingThreshold := now.Add(24 * time.Hour)

		sqlStmt := `SELECT id, description, status, due_date FROM tasks WHERE status != ? AND due_date IS NOT NULL ORDER BY due_date ASC`
		statusDone := "done"

		rows, err := db.DB.Query(sqlStmt, statusDone)
		if err != nil {
			log.Fatalf("Failed to query tasks for checking: %v", err)
		}
		defer rows.Close()

		foundTasks := 0
		overdueTasksOutput := []string{}
		upcomingTasksOutput := []string{}

		for rows.Next() {
			var id int
			var description, status, dueDateStr string

			err := rows.Scan(&id, &description, &status, &dueDateStr)
			if err != nil {
				log.Printf("Warning: Failed to scan row during overdue check: %v", err)
				continue
			}

			dueDateTime, parseErr := time.Parse("2006-01-02 15:04:05", dueDateStr)
			if parseErr != nil {
				log.Printf("Warning: Could not parse due date '%s' from DB for task ID %d: %v", dueDateStr, id, parseErr)
				continue
			}

			if dueDateTime.Before(now) {
				foundTasks++
				line := fmt.Sprintf("  - ID %d: %s (Due: %s)", id, description, dueDateTime.Format("Mon, 02 Jan - 15:04"))
				overdueTasksOutput = append(overdueTasksOutput, line)
			}
		}
		if err = rows.Err(); err != nil {
			log.Fatalf("Error iterating over tasks (overdue check): %v", err)
		}
		rows.Close()

		rows, err = db.DB.Query(sqlStmt, statusDone)
		if err != nil {
			log.Fatalf("Failed to re-query tasks for upcoming check: %v", err)
		}

		for rows.Next() {
			var id int
			var description, status, dueDateStr string

			err := rows.Scan(&id, &description, &status, &dueDateStr)
			if err != nil {
				log.Printf("Warning: Failed to scan row during upcoming check: %v", err)
				continue
			}

			dueDateTime, parseErr := time.Parse("2006-01-02 15:04:05", dueDateStr)
			if parseErr != nil {
				continue
			}

			if !dueDateTime.Before(now) && dueDateTime.Before(upcomingThreshold) {
				foundTasks++
				line := fmt.Sprintf("  - ID %d: %s (Due: %s)", id, description, dueDateTime.Format("Mon, 02 Jan - 15:04"))
				upcomingTasksOutput = append(upcomingTasksOutput, line)
			}
		}
		if err = rows.Err(); err != nil {
			log.Fatalf("Error iterating over tasks (upcoming check): %v", err)
		}

		if foundTasks > 0 {
			fmt.Println("---------------")

			fmt.Println("🚨 Overdue Tasks:")
			if len(overdueTasksOutput) > 0 {
				fmt.Println(strings.Join(overdueTasksOutput, "\n"))
			} else {
				fmt.Println("  (None)")
			}

			fmt.Println("\n✨ Upcoming Tasks (due within 24 hours):")
			if len(upcomingTasksOutput) > 0 {
				fmt.Println(strings.Join(upcomingTasksOutput, "\n"))
			} else {
				fmt.Println("  (None)")
			}
			fmt.Println("---------------")

		} else {
			if !silentIfClear {
				fmt.Println("👍 No overdue or upcoming tasks found. You're all clear!")
			}
		}

		silentIfClear = false
	},
}

func init() {
	AddCommand(taskCmd)

	taskCmd.AddCommand(taskAddCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskDoneCmd)
	taskCmd.AddCommand(taskRemoveCmd)
	taskCmd.AddCommand(taskCheckCmd)

	taskAddCmd.Flags().StringVarP(&dueDateStr, "due", "d", "", "Due date: 'YYYY-MM-DD', 'YYYY-MM-DD HH:MM:SS', or relative (e.g., '1d', '2h30m', '1w')")
	taskCheckCmd.Flags().BoolVarP(&silentIfClear, "silent-if-clear", "s", false, "Suppress output if no overdue or upcoming tasks are found")
}
