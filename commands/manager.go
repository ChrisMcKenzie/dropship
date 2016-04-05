// Copyright (c) 2016 "ChrisMcKenzie"
// This file is part of Dropship.
//
// Dropship is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License v3 as
// published by the Free Software Foundation
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
package commands

import (
	"fmt"

	"github.com/ChrisMcKenzie/dropship/manager"
	"github.com/spf13/cobra"
)

var StoreURL string

var managerCmd = &cobra.Command{
	Use:   "manage",
	Short: "starts a manager api and interface",
	Run:   managerC,
}

func init() {
	managerCmd.PersistentFlags().StringVar(&StoreURL, "store", "", "where dropship stores its data (ie. consul://127.0.0.1:8500)")
}

func managerC(c *cobra.Command, args []string) {
	printInitSummary()
	manager.Start(StoreURL)
}

func printInitSummary() {
	fmt.Printf(
		"Dropship\n  version: %s\n  build date: %s\nSettings\n  storage url: %s\n",
		version,
		buildDate,
		StoreURL,
	)
}
