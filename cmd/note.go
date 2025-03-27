package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/57ajay/goTask/db"
	"github.com/spf13/cobra"
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Manage your notes",
	Long:  `Add, list, and remove quick notes.`,
	Run:   noteListCmd.Run,
}

var noteAddCmd = &cobra.Command{
	Use:   "add [note content]",
	Short: "Add a new note",
	Long:  `Add a new note. The content should be provided as arguments.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		content := strings.Join(args, " ")

		sqlStmt := `INSERT INTO notes (content) VALUES (?)`
		_, err := db.DB.Exec(sqlStmt, content)
		if err != nil {
			log.Fatalf("Failed to add note: %v", err)
		}
		fmt.Printf("Added note: \"%s\"\n", content)
	},
}

var noteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your notes",
	Long:  `Lists all notes currently stored.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Your Notes:")
		fmt.Println("-----------")

		sqlStmt := `SELECT id, content, created_at FROM notes ORDER BY id ASC`
		rows, err := db.DB.Query(sqlStmt)
		if err != nil {
			log.Fatalf("Failed to list notes: %v", err)
		}
		defer rows.Close()

		noteCount := 0
		for rows.Next() {
			noteCount++
			var id int
			var content, createdAt string

			err := rows.Scan(&id, &content, &createdAt)
			if err != nil {
				log.Printf("Warning: Failed to scan row: %v", err)
				continue
			}
			fmt.Printf("%d. %s (Added: %s)\n", id, content, createdAt)
		}

		if err = rows.Err(); err != nil {
			log.Fatalf("Error iterating over notes: %v", err)
		}

		if noteCount == 0 {
			fmt.Println("You have no notes yet!")
		}
		fmt.Println("-----------")
	},
}

var noteRemoveCmd = &cobra.Command{
	Use:   "remove [note_id]",
	Short: "Remove a note",
	Long:  `Remove a note permanently by providing its ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		idStr := args[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatalf("Invalid note ID provided: '%s'. Please provide a number.", idStr)
		}

		sqlStmt := `DELETE FROM notes WHERE id = ?`
		result, err := db.DB.Exec(sqlStmt, id)
		if err != nil {
			log.Fatalf("Failed to remove note %d: %v", id, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("Warning: Could not check rows affected after delete: %v", err)
			fmt.Printf("Attempted to remove note %d.\n", id)
		} else if rowsAffected == 0 {
			fmt.Printf("Note with ID %d not found.\n", id)
		} else {
			fmt.Printf("Removed note %d.\n", id)
		}
	},
}

func init() {
	AddCommand(noteCmd)

	noteCmd.AddCommand(noteAddCmd)
	noteCmd.AddCommand(noteListCmd)
	noteCmd.AddCommand(noteRemoveCmd)
}
