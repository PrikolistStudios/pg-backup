package cmd

import (
	"fmt"
	"os"

	"github.com/PrikolistStudios/pg-backup/internal/app"
	"github.com/spf13/cobra"
)

var config = app.NewConfig()
var mode = "backup"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pg-backup [arguments...]",
	Short: "Removal and backup of PostgreSQL databases",
	Long: `This CLI tool is used for easy removal and backup of PostgreSQL databases. 
It supports action on multiple databases in one command and globbing

Each argument should be the database name, and the user must have necessary rights to perform an action on it 
(read to backup, drop to remove). Argument can also be a glob string to remove many databases.

When removing the database, it is not backed up implicitly. Database backups are stored in current working directory.
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if mode != "backup" && mode != "remove" {
			_, _ = fmt.Fprintf(os.Stderr, "Error: --mode must be either 'backup' or 'remove'\n")
			return
		}

		var err error
		if mode == "backup" {
			err = app.BackupDatabases(args, config)
		} else {
			err = app.RemoveDatabases(args, config)
		}

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, `Error: %v

Errors were encountered. Check the log for details`, err)
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
