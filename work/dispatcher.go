package work

import (
	"errors"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/ChrisMcKenzie/dropship/hook"
	"github.com/ChrisMcKenzie/dropship/installer"
	"github.com/ChrisMcKenzie/dropship/lock"
	"github.com/ChrisMcKenzie/dropship/service"
	"github.com/ChrisMcKenzie/dropship/updater"
	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

// Dispatcher is responsible for managing a given services state and
// sending work to the Runner pool
type Dispatcher struct {
	config     service.Config
	ticker     *time.Ticker
	task       *Runner
	hash       string
	duration   time.Duration
	wg         *sync.WaitGroup
	shutdownCh <-chan struct{}
}

func NewDispatcher(cfg service.Config, t *Runner, wg *sync.WaitGroup, shutdownCh <-chan struct{}) (*Dispatcher, error) {
	w := Dispatcher{
		config:     cfg,
		task:       t,
		shutdownCh: shutdownCh,
		wg:         wg,
	}

	var err error
	w.duration, err = time.ParseDuration(cfg.CheckInterval)
	if err != nil {
		return nil, err
	}

	go w.start()

	return &w, nil
}

func (w *Dispatcher) start() {
	for {
		select {
		case _, ok := <-w.shutdownCh:
			if !ok {
				log.Printf("Shutting down dispatcher for %s", w.config.Name)
				w.wg.Done()
				return
			}
		case <-time.After(w.duration):
			w.task.Do(w)
		}
	}
}

func (w *Dispatcher) Work() {
	log.Printf("[INF]: Starting Update check for %s...", w.config.Name)
	user := viper.GetString("rackspaceUser")
	key := viper.GetString("rackspaceKey")
	region := viper.GetString("rackspaceRegion")

	u := updater.NewRackspaceUpdater(user, key, region)
	opts := &updater.Options{w.config.Artifact[0].Bucket, w.config.Artifact[0].Path}

	isOutOfDate, err := u.IsOutdated(w.config.Hash, opts)
	if err != nil {
		log.Printf("[ERR]: Unable to check updates for %s %v", w.config.Name, err)
		return
	}

	if isOutOfDate {
		if w.config.Sequential {
			log.Printf("[INF]: Acquiring lock for %s", w.config.Name)
			l, err := lock.NewConsulLocker(w.config.Name, api.DefaultConfig())
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

		log.Printf("[INF]: Downloading update for %s...", w.config.Name)
		fr, meta, err := u.Download(opts)
		if err != nil {
			log.Printf("[ERR]: Unable to download update for %s %v", w.config.Name, err)
			return
		}

		if w.config.PreCommand != "" {
			log.Printf("[INF]: preCommand has been deprecated.")
			res, err := executeCommand(w.config.PreCommand)
			if err != nil {
				log.Printf("[ERR]: Unable to execute preCommand. %v", err)
			}
			log.Printf("[INF]: preCommand executed successfully. %v", res)
		}

		err = runHooks(w.config.BeforeHooks, w.config)
		if err != nil {
			log.Printf("[ERR]: Unable to execute beforeHooks. %v", err)
		}

		i, err := getInstaller(meta.ContentType)
		if err != nil {
			log.Printf("[ERR]: %s for %s", w.config.Name, err)
			return
		}

		filesWritten, err := i.Install(w.config.Artifact[0].Destination, fr)
		if err != nil {
			log.Printf("[ERR]: Unable to install update for %s %s", w.config.Name, err)
		}

		if w.config.PostCommand != "" {
			log.Printf("[INF]: postCommand has been deprecated.")
			defer func() {
				res, err := executeCommand(w.config.PostCommand)
				if err != nil {
					log.Printf("[ERR]: Unable to execute postCommand. %v", err)
				}
				log.Printf("[INF]: postCommand executed successfully. %v", res)
			}()
		}

		log.Printf("[INF]: Update for %s installed successfully. [hash: %s] [files written: %d]", w.config.Name, meta.Hash, filesWritten)
		// TODO(ChrisMcKenzie): hashes should be stored somewhere more
		// permanent.
		w.config.Hash = meta.Hash

		err = runHooks(w.config.AfterHooks, w.config)
		if err != nil {
			log.Printf("[ERR]: Unable to execute beforeHooks. %v", err)
		}
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
	case "application/x-gzip", "application/octet-stream":
		var installer installer.TarInstaller
		return installer, nil
	}

	return nil, errors.New("Unable to determine installation method from file type")
}

func runHooks(hooks []service.Hook, service service.Config) error {
	for _, h := range hooks {
		for hookName, config := range h {
			hook := hook.GetHookByName(hookName)
			if hook != nil {
				log.Printf("[INF]: Executing \"%s\" hook with %+v", hookName, config)
				hook.Execute(config, service)
			}
		}
	}
	return nil
}
