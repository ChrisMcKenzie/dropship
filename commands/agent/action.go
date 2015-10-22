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

var dataDir, rackspaceUser, rackspaceKey, rackspaceRegion string

func init() {
	Command.Flags().StringVarP(&dataDir, "data-dir", "d", "./", "directory that agent will look at for service configs")
	Command.Flags().StringVar(&rackspaceUser, "rackspace-user", "", "user to use for rackspace repo")
	Command.Flags().StringVar(&rackspaceKey, "rackspace-key", "", "key to use for rackspace repo")
	Command.Flags().StringVar(&rackspaceRegion, "rackspace-region", "IAD", "region to use for rackspace repo")
}

func Action(cmd *cobra.Command, args []string) {
	services := loadServices(dataDir)

	rackspace.Setup(
		rackspaceUser,
		rackspaceKey,
		rackspaceRegion,
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
