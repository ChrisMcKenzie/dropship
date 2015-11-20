package agent

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/ChrisMcKenzie/dropship/dropship"
)

// Dispatcher is responsible for managing a given services state and
// sending work to the Runner pool
type Dispatcher struct {
	config     dropship.Config
	task       *Runner
	hash       string
	duration   time.Duration
	wg         *sync.WaitGroup
	shutdownCh <-chan struct{}
}

func NewDispatcher(cfg dropship.Config, t *Runner, wg *sync.WaitGroup, shutdownCh <-chan struct{}) (*Dispatcher, error) {
	w := Dispatcher{
		config:     cfg,
		task:       t,
		shutdownCh: shutdownCh,
		wg:         wg,
	}

	var err error
	w.duration, err = time.ParseDuration(cfg.CheckInterval)
	if err != nil {
		return nil, fmt.Errorf("Dispatcher: Failed to start %s", err)
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

	u := w.config.Updater

	isOutOfDate, err := u.IsOutdated(w.config.Hash, w.config.Artifact)
	if err != nil {
		log.Printf("[ERR]: Unable to check updates for %s %v", w.config.Name, err)
		return
	}

	if isOutOfDate {
		if w.config.Sequential {
			log.Printf("[INF]: Acquiring lock for %s", w.config.Name)
			l := w.config.Locker
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
		fr, meta, err := u.Download(w.config.Artifact)
		if err != nil {
			log.Printf("[ERR]: Unable to download update for %s %v", w.config.Name, err)
			return
		}

		// Deprecated
		if w.config.PreCommand != "" {
			log.Printf("[WARN]: preCommand has been deprecated.")
			res, err := executeCommand(w.config.PreCommand)
			if err != nil {
				log.Printf("[ERR]: Unable to execute preCommand. %v", err)
			} else {
				log.Printf("[INF]: preCommand executed successfully. %v", res)
			}
		}

		err = runHooks(w.config.BeforeHooks, w.config)
		if err != nil {
			log.Printf("[ERR]: Unable to execute beforeHooks. %v", err)
		}

		contentType := meta.ContentType
		if ct, ok := w.config.Artifact["content-type"]; ok {
			contentType = ct
		}

		i, err := getInstaller(contentType)
		if err != nil {
			log.Printf("[ERR]: %s for %s", w.config.Name, err)
			return
		}

		filesWritten, err := i.Install(w.config.Artifact["destination"], fr)
		if err != nil {
			log.Printf("[ERR]: Unable to install update for %s %s", w.config.Name, err)
		}

		// Deprecated
		if w.config.PostCommand != "" {
			log.Printf("[WARN]: postCommand has been deprecated.")
			defer func() {
				res, err := executeCommand(w.config.PostCommand)
				if err != nil {
					log.Printf("[ERR]: Unable to execute postCommand. %v", err)
				} else {
					log.Printf("[INF]: postCommand executed successfully. %v", res)
				}
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

// Deprecated
func executeCommand(c string) (string, error) {
	cmd := strings.Fields(c)
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	return string(out), err
}

func getInstaller(contentType string) (dropship.Installer, error) {
	switch contentType {
	case "application/x-gzip", "application/octet-stream":
		var installer dropship.TarInstaller
		return installer, nil
	default:
		var installer dropship.FileInstaller
		return installer, nil
	}

	return nil, errors.New("Unable to determine installation method from file type")
}

func runHooks(hooks []dropship.HookDefinition, service dropship.Config) error {
	for _, h := range hooks {
		for hookName, config := range h {
			hook := dropship.GetHookByName(hookName)
			if hook != nil {
				log.Printf("[INF]: Executing \"%s\" hook with %+v", hookName, config)
				err := hook.Execute(config, service)
				if err != nil {
					log.Printf("[ERR]: Unable to execute \"%s\" hook %v", hookName, err)
				}
			}
		}
	}
	return nil
}
