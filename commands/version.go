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
