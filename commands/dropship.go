package commands

import (
	"github.com/ChrisMcKenzie/dropship/commands/agent"
	"github.com/spf13/cobra"
)

var DropshipCmd = &cobra.Command{
	Use:   "dropship",
	Short: "dropship deploys your code",
	Long:  "dropship monitors and automatically updates your code",
}

func Execute() {
	AddCommands()
	DropshipCmd.Execute()
}

func AddCommands() {
	DropshipCmd.AddCommand(agent.Command)
}
