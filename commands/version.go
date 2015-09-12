package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var timeLayout string // the layout for time.Time

var version = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Dropship",
	Long:  `All software has versions. This is Dropship's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("version: 1")
	},
}
