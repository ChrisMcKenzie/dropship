package commands

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/ChrisMcKenzie/dropship/commands/agent"
	"github.com/ChrisMcKenzie/dropship/dropship"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "starts automatic checks and update",
	Run:   agentC,
}

func agentC(c *cobra.Command, args []string) {
	InitializeConfig()
	root := viper.GetString("servicePath")
	services, err := dropship.LoadServices(root)
	if err != nil {
		log.Fatalln(err)
	}

	t := agent.NewRunner(len(services))
	shutdownCh := make(chan struct{})

	var wg sync.WaitGroup
	wg.Add(len(services))
	for _, s := range services {
		_, err := agent.NewDispatcher(s, t, &wg, shutdownCh)
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
