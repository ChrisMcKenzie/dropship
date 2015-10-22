package agent

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ChrisMcKenzie/dropship/repo/rackspace"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "agent",
	Short: "Start updater service",
	Long: `Starts an agent that will check and update all defined services
on a given interval
	`,
	Run: Action,
}

var configPath string

func init() {
	Command.Flags().StringVar(&configPath, "config", "", "path to config file")
}

func Action(cmd *cobra.Command, args []string) {
	if configPath == "" {
		log.Fatal("config must be supplied")
	}

	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	services := loadServices(config.ConfigDir)

	rackspace.Setup(
		config.Rackspace.User,
		config.Rackspace.Key,
		config.Rackspace.Region,
	)

	ch := make(chan struct{})
	var wg *sync.WaitGroup = &sync.WaitGroup{}
	wg.Add(1)

	for _, service := range services {
		log.Println("setting up service", service.Id)
		updater, err := setup(service, ch)
		if err != nil {
			panic(err)
		}
		go updater.Start(wg)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func(ch chan struct{}, sigs chan os.Signal) {
		<-sigs
		close(ch)
		wg.Done()
	}(ch, sigs)

	wg.Wait()
}
