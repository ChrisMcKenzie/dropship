package commands

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/ChrisMcKenzie/dropship/commands/agent"
	"github.com/ChrisMcKenzie/dropship/dropship"
	"github.com/spf13/cobra"
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

	if cfg.Rackspace != nil {
		log.Println("[WARN]: The Rackspace config item has been deprecated and will be removed in future versions. please use the repo directive. ")
		updaters["rackspace"] = dropship.NewRackspaceUpdater(cfg.Rackspace)
	}
	initializeUpdaters(cfg.Repos)
	initializeLockers(cfg.Locks)

	services, err := dropship.LoadServices(cfg.ServicePath)
	if err != nil {
		log.Fatalln(err)
	}

	runner := agent.NewRunner(len(services))
	shutdownCh := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(len(services))
	for _, service := range services {
		log.Printf("[INF]: Starting updater for %s", service.Name)
		var ok bool
		service.Updater, ok = updaters[service.Artifact["type"]]
		if !ok {
			log.Fatalf("[ERR]: Unable to find updater %s", service.Artifact["type"])
		}

		// Try and use consul config but nothing exists use default consul config.
		//
		// TODO(ChrisMcKenzie): this is ugly and should support more than just consul
		// lockers
		service.Locker, ok = lockers["consul"]
		if !ok {
			var err error
			service.Locker, err = dropship.NewConsulLocker(map[string]string{})
			if err != nil {
				log.Fatalf("[ERR]: Unable to initilize locker: %s", err)
			}
		}

		_, err := agent.NewDispatcher(service, runner, &wg, shutdownCh)
		if err != nil {
			log.Fatal(err)
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
