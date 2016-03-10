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
	"io/ioutil"
	"log"
	"os"

	"github.com/hashicorp/hcl"
	"github.com/spf13/cobra"
)

type RepoConfig map[string]string
type LockConfig map[string]string

type Config struct {
	ManagerURL  string                `hcl:"manager_url"`
	ServicePath string                `hcl:"service_path"`
	Rackspace   map[string]string     `hcl:"rackspace"`
	Repos       map[string]RepoConfig `hcl:"repo"`
	Locks       map[string]LockConfig `hcl:"lock"`
}

var DropshipCmd = &cobra.Command{
	Use:   "dropship",
	Short: "A tool for automated and distributed artifact deployment",
	Long: `

Dropship is a distrubuted deployment system and interface
allowing users to download build artifacts from any sort of cdn and 
install it on any number of hosts in a distributed way.`,
}

var CfgFile string

func init() {
	DropshipCmd.PersistentFlags().StringVar(&CfgFile, "config", "/etc/dropship.d/dropship.hcl", "config file")
}

func Execute() {
	AddCommands()
	if err := DropshipCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func AddCommands() {
	DropshipCmd.AddCommand(agentCmd)
	DropshipCmd.AddCommand(versionCmd)
	DropshipCmd.AddCommand(managerCmd)
}

func InitializeConfig() *Config {
	var cfg Config
	cfgData, err := ioutil.ReadFile(CfgFile)
	if err != nil {
		log.Fatalln("Unable to locate Config File. make sure you specify it using the --config flag")
		return nil
	}
	err = hcl.Decode(&cfg, string(cfgData))

	if err != nil {
		log.Fatalln("Unable to parse Config File.")
		return nil
	}

	return &cfg
}
