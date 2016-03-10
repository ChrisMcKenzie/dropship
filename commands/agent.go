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
	"os"
	"os/signal"
	"sync"

	"github.com/ChrisMcKenzie/dropship/agent"
	"github.com/ChrisMcKenzie/dropship/dropship"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var updaters map[string]dropship.Updater = make(map[string]dropship.Updater)
var lockers map[string]dropship.Locker = make(map[string]dropship.Locker)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "starts automatic checks and update",
	Run:   agentC,
}

func agentC(c *cobra.Command, args []string) {
	cfg := InitializeConfig()

	conn, err := grpc.Dial(cfg.ManagerURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := dropship.NewRpcServiceClient(conn)

	if cfg.Rackspace != nil {
		log.Warn("The Rackspace config item has been deprecated and will be removed in future versions. please use the repo directive. ")
		updaters["rackspace"] = dropship.NewRackspaceUpdater(cfg.Rackspace)
	}
	initializeUpdaters(cfg.Repos)
	initializeLockers(cfg.Locks)

	services, err := dropship.LoadServices(cfg.ServicePath)
	if err != nil {
		log.Fatal(err)
	}

	runner := agent.NewRunner(len(services))
	shutdownCh := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(len(services))
	for _, service := range services {
		log.Infof("Starting updater for %s", service.Name)
		var ok bool
		service.Updater, ok = updaters[service.Artifact["type"]]
		if !ok {
			log.Errorf("Unable to find updater %s", service.Artifact["type"])
		}

		_, err := client.RegisterService(context.Background(), &dropship.Service{service.Name})
		if err != nil {
			log.Error(err)
		}

		// Try and use consul config but nothing exists use default consul config.
		//
		// TODO(ChrisMcKenzie): this is ugly and should support more than just consul
		// lockers
		service.Locker, ok = lockers["consul"]
		if !ok {
			var err error
			service.Locker, err = dropship.NewConsulLocker(cfg.Locks["consul"])
			if err != nil {
				log.Errorf("Unable to initilize locker: %s", err)
			}
		}

		_, err = agent.NewDispatcher(service, runner, &wg, shutdownCh)
		if err != nil {
			log.Error(err)
		}
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs
	close(shutdownCh)
	wg.Wait()

	runner.Shutdown()
}

func initializeLockers(configs map[string]LockConfig) {
	for name, value := range configs {
		switch name {
		case "consul":
			consulLock, err := dropship.NewConsulLocker(value)
			if err != nil {
				log.Fatalf("[ERR]: Unable to initialize Consul Locker: %s", err)
			}
			lockers[name] = consulLock
		}
	}
}

func initializeUpdaters(configs map[string]RepoConfig) {
	for name, value := range configs {
		switch name {
		case "rackspace":
			updaters[name] = dropship.NewRackspaceUpdater(value)
		case "s3":
			updaters[name] = dropship.NewS3Updater(value)
		}
	}
}
