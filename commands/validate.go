package commands

import (
	"fmt"
	"io/ioutil"

	"github.com/ChrisMcKenzie/dropship/couriers"
	"github.com/spf13/cobra"
)

var validate = &cobra.Command{
	Use:   "validate",
	Short: "Parse and validate the given dropship.yml",
	Long:  `Parse and validate the given dropship.yml`,
	Run: func(cmd *cobra.Command, args []string) {
		bytes, err := ioutil.ReadFile("dropship.yml")
		d, err := couriers.ParseDeployment(bytes)
		if err != nil {
			fmt.Println(err)
		}

		for key, _ := range d.Servers {
			fmt.Printf(
				"Environment: %s\n",
				key,
			)
		}

		fmt.Println("")

		for _, val := range d.Commands {
			fmt.Printf(
				"Command: %s\n",
				val,
			)
		}
	},
}
