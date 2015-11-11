package commands

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/ChrisMcKenzie/dropship/service"
	"github.com/ChrisMcKenzie/dropship/work"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "starts automatic checks and update",
	Run:   agent,
}

func agent(c *cobra.Command, args []string) {
	InitializeConfig()
	root := viper.GetString("servicePath")
	services, err := service.LoadServices(root)
	if err != nil {
		log.Fatalln(err)
	}

	t := work.NewRunner(len(services))
	shutdownCh := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(len(services))
	for _, s := range services {
		_, err := work.NewDispatcher(s, t, &wg, shutdownCh)
		if err != nil {
			log.Fatal(err)
		}
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs
	close(shutdownCh)
	wg.Wait()

	t.Shutdown()
}
