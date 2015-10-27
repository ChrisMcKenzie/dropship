package commands

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/ChrisMcKenzie/dropship/service"
	"github.com/hashicorp/hcl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const maxGoRoutines = 10

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "starts automatic checks and update",
	Run:   agent,
}

func agent(c *cobra.Command, args []string) {
	InitializeConfig()

	root := viper.GetString("servicePath")
	services, err := loadServices(root)
	if err != nil {
		log.Fatalln(err)
	}

	t := service.NewRunner(len(services))
	shutdownCh := make(chan struct{})

	for _, s := range services {
		service.NewDispatcher(s, t, shutdownCh)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs
	close(shutdownCh)

	t.Shutdown()
}

func loadServices(root string) (d []service.Config, err error) {
	files, _ := filepath.Glob(root + "/*.hcl")
	for _, file := range files {
		var data []byte
		data, err = ioutil.ReadFile(file)
		if err != nil {
			return
		}

		var deploy struct {
			Services []service.Config `hcl:"service,expand"`
		}
		hcl.Decode(&deploy, string(data))
		d = append(d, deploy.Services...)
	}
	return
}
