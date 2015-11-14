package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version    string
	buildDate  string
	commitHash string
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Dropship Version: %s\nBuild Date: %s\nCommit: %s\n", version, buildDate, commitHash)
	},
}
