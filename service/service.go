package service

import (
	"errors"
	"log"
	"time"

	"github.com/ChrisMcKenzie/dropship/installer"
	"github.com/ChrisMcKenzie/dropship/updater"
	"github.com/spf13/viper"
)

type Artifact struct {
	Type        string `hcl:",key"`
	Bucket      string `hcl:"bucket"`
	Path        string `hcl:"path"`
	Destination string `hcl:"destination"`
}

type Config struct {
	Name          string   `hcl:",key"`
	CheckInterval string   `hcl:"checkInterval"`
	Artifact      Artifact `hcl:"artifact,expand"`
}

type Dispatcher struct {
	config     Config
	ticker     *time.Ticker
	task       *Runner
	hash       string
	shutdownCh <-chan struct{}
}

func NewDispatcher(cfg Config, t *Runner, shutdownCh <-chan struct{}) (*Dispatcher, error) {
	w := Dispatcher{
		config:     cfg,
		task:       t,
		shutdownCh: shutdownCh,
	}

	dur, err := time.ParseDuration(cfg.CheckInterval)
	if err != nil {
		return nil, err
	}
	w.ticker = time.NewTicker(dur)

	go w.start()

	return &w, nil
}

func (w Dispatcher) start() {
	for {
		select {
		case <-w.ticker.C:
			w.task.Do(w)
		case _, ok := <-w.shutdownCh:
			if !ok {
				log.Printf("Shutting down dispatcher for %s", w.config.Name)
				w.ticker.Stop()
				return
			}
		}
	}
}

func (w Dispatcher) Work() {
	log.Printf("[INF]: Starting Update check for %s...", w.config.Name)
	user := viper.GetString("rackspaceUser")
	key := viper.GetString("rackspaceKey")
	region := viper.GetString("rackspaceRegion")

	u := updater.NewRackspaceUpdater(user, key, region)
	opts := &updater.Options{w.config.Artifact.Bucket, w.config.Artifact.Path}

	isOutOfDate, err := u.IsOutdated(w.hash, opts)
	if err != nil {
		log.Printf("[ERR]: Unable to check updates for %s %v", w.config.Name, err)
		return
	}

	if isOutOfDate {
		log.Printf("[INF]: Installing update for %s...", w.config.Name)
		fr, meta, err := u.Download(opts)
		if err != nil {
			log.Printf("[ERR]: Unable to download update for %s %v", w.config.Name, err)
			return
		}

		i, err := getInstaller(meta.ContentType)
		if err != nil {
			log.Printf("[ERR]: %s for %s", w.config.Name, err)
			return
		}

		filesWritten, err := i.Install(w.config.Artifact.Destination, fr)
		if err != nil {
			log.Printf("[ERR]: Unable to install update for %s %s", w.config.Name, err)
		}

		log.Printf("[INF]: Update for %s installed successfully. [files written: %d]", w.config.Name, filesWritten)
		w.hash = meta.Hash
	}
}

func getInstaller(contentType string) (installer.Installer, error) {
	switch contentType {
	case "application/x-gzip":
		var installer installer.TarInstaller
		return installer, nil
	}

	return nil, errors.New("Unable to determine installation method from file type")
}
