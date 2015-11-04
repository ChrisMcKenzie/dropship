package service

import (
	"errors"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/ChrisMcKenzie/dropship/installer"
	"github.com/ChrisMcKenzie/dropship/lock"
	"github.com/ChrisMcKenzie/dropship/updater"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

// Dispatcher is responsible for managing a given services state and
// sending work to the Runner pool
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

func (w *Dispatcher) start() {
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

func (w *Dispatcher) Work() {
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
		if w.config.Sequential {
			log.Printf("[INF]: Acquiring lock for %s", w.config.Name)
			l, err := lock.NewConsulLocker("dropship/services/"+w.config.Name, api.DefaultConfig())
			if err != nil {
				log.Printf("[ERR]: Unable to retreive update lock. %v", err)
				return
			}
			_, err = l.Acquire(w.shutdownCh)
			if err != nil {
				log.Printf("[ERR]: Unable to retreive update lock. %v", err)
				return
			}
			defer l.Release()
		}

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

		if w.config.PostCommand != "" {
			res, err := executeCommand(w.config.PostCommand)
			if err != nil {
				log.Printf("[ERR]: Unable to execute postComment. %v", err)
			}
			log.Printf("[INF]: postCommand executed successfully. %v", res)
		}

		log.Printf("[INF]: Update for %s installed successfully. [hash: %s] [files written: %d]", w.config.Name, meta.Hash, filesWritten)
		w.hash = meta.Hash
	} else {
		log.Printf("[INF]: %s is up to date", w.config.Name)
	}
}

func executeCommand(c string) (string, error) {
	cmd := strings.Fields(c)
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	return string(out), err
}

func getInstaller(contentType string) (installer.Installer, error) {
	switch contentType {
	case "application/x-gzip":
		var installer installer.TarInstaller
		return installer, nil
	}

	return nil, errors.New("Unable to determine installation method from file type")
}
