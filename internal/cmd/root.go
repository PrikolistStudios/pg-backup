package cmd

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/PrikolistStudios/pg-backup/internal/app"
	"github.com/spf13/cobra"
)

//go:embed long-description.txt
var longDescription string
var config = app.NewConfig()
var mode = "backup"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pg-backup [arguments...]",
	Short: "Removal and backup of PostgreSQL databases",
	Long:  longDescription,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if mode != "backup" && mode != "remove" {
			_, _ = fmt.Fprintf(os.Stderr, "Error: --mode must be either 'backup' or 'remove'\n")
			return
		}

		// Connect to db.
		conn, err := app.CreateConnection(config)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: failed to connect to database. Check connection parameters.\n")
			return
		}
		defer func(conn *sql.DB) {
			_ = conn.Close()
		}(conn)

		// Get valid database names.
		names, err := app.FilterPatterns(args, conn)
		var globErr app.ErrAccumulatedErrors
		if errors.As(err, &globErr) {
			_, _ = fmt.Fprintf(os.Stderr, "Error occurred while resolving glob patterns:\n%s", globErr)
		} else if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error occurred while fetching databases list: %s", err)
			return
		}

		if len(names) == 0 {
			fmt.Println("No databases found. Aborting.")
			return
		}

		// List databases and prompt the user to confirm action.
		var confirm bool
		prompt := &survey.Confirm{
			Message: fmt.Sprintf("Perform the action (%s) on the following databases?\n%s\n", mode, strings.Join(names, "\n")),
			Default: false,
		}

		err = survey.AskOne(prompt, &confirm)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Error confirming action. Aborting.")
			return
		} else if !confirm {
			fmt.Println("Aborting.")
			return
		}

		fmt.Printf("Performing %s\n", mode)

		// Finally perform the action.
		var action app.DatabaseAction
		if mode == "backup" {
			action = app.NewBackupAction(config)
		} else {
			action = app.NewRemoveAction(config.ForceRemove, conn)
		}

		err = app.PerformDatabasesAction(names, action)

		// List databases with unsuccessful action.
		var dbErr app.ErrAccumulatedErrors
		if errors.As(err, &dbErr) {
			_, _ = fmt.Fprintf(os.Stderr, "Error occurred while performing action on these databases:\n%s\n\n Check access rights and database existence.", globErr)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&config.Host, "host", "H", "localhost", "Database server host")
	rootCmd.Flags().StringVarP(&config.Port, "port", "p", "5432", "Database server port")
	rootCmd.Flags().StringVarP(&config.User, "user", "U", "postgres", "Database user")
	rootCmd.Flags().StringVarP(&config.Password, "password", "P", "postgres", "User password")
	rootCmd.Flags().StringVarP(&config.Database, "database", "d", "postgres", "Database to connect to while performing removal")
	rootCmd.Flags().StringVarP(&mode, "mode", "m", "backup",
		`Action performed on database. Can be either 'backup' of 'remove'`)
	rootCmd.Flags().BoolVarP(&config.ForceRemove, "force", "f", false, "Remove database even if it has active connections")
}
