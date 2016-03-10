// Copyright (c) 2016 "ChrisMcKenzie"
// This file is part of Dropship.
//
// Dropship is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License v3 as
// published by the Free Software Foundation
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
package agent

import (
	"errors"
	"fmt"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

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
				log.Infof("Shutting down dispatcher for %s", w.config.Name)
				w.wg.Done()
				return
			}
		case <-time.After(w.duration):
			w.task.Do(w)
		}
	}
}

func (w *Dispatcher) Work() {
	log.Infof("Starting Update check for %s...", w.config.Name)

	u := w.config.Updater

	isOutOfDate, err := u.IsOutdated(w.config.Hash, w.config.Artifact)
	if err != nil {
		log.Errorf("Unable to check updates for %s %v", w.config.Name, err)
		return
	}

	if isOutOfDate {
		if w.config.Sequential {
			log.Infof("Acquiring lock for %s", w.config.Name)
			l := w.config.Locker
			if err != nil {
				log.Errorf("Unable to retreive update lock. %v", err)
				return
			}
			_, err = l.Acquire(w.shutdownCh)
			if err != nil {
				log.Errorf("Unable to retreive update lock. %v", err)
				return
			}
			defer l.Release()
		}

		log.Infof("Downloading update for %s...", w.config.Name)
		fr, meta, err := u.Download(w.config.Artifact)
		if err != nil {
			log.Errorf("Unable to download update for %s %v", w.config.Name, err)
			return
		}

		// Deprecated
		err = runHooks(w.config.BeforeHooks, w.config)
		if err != nil {
			log.Errorf("Unable to execute beforeHooks. %v", err)
		}

		contentType := meta.ContentType
		if ct, ok := w.config.Artifact["content-type"]; ok {
			contentType = ct
		}

		i, err := getInstaller(contentType)
		if err != nil {
			log.Errorf("%s for %s", w.config.Name, err)
			return
		}

		filesWritten, err := i.Install(w.config.Artifact["destination"], fr)
		if err != nil {
			log.Errorf("Unable to install update for %s %s", w.config.Name, err)
		}

		log.Infof("Update for %s installed successfully. [hash: %s] [files written: %d]", w.config.Name, meta.Hash, filesWritten)
		// TODO(ChrisMcKenzie): hashes should be stored somewhere more
		// permanent.
		// This should be sent to the manager
		w.config.Hash = meta.Hash

		err = runHooks(w.config.AfterHooks, w.config)
		if err != nil {
			log.Errorf("Unable to execute beforeHooks. %v", err)
		}

		if w.config.UpdateTTL != "" {
			log.Infof("Waiting %s before releasing lock and allowing next deployment.", w.config.UpdateTTL)
			ttl, err := time.ParseDuration(w.config.UpdateTTL)
			if err != nil {
				log.Errorf("Failed to parse updateTTL make sure it is a valid duration in seconds")
			}
			<-time.After(ttl)
		}
	} else {
		log.Infof("%s is up to date", w.config.Name)
	}
}

func getInstaller(contentType string) (dropship.Installer, error) {
	switch contentType {
	case "application/x-gzip", "application/octet-stream", "application/gzip":
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
				log.Infof("Executing \"%s\" hook with %+v", hookName, config)
				err := hook.Execute(config, service)
				if err != nil {
					log.Error("Unable to execute \"%s\" hook %v", hookName, err)
				}
			}
		}
	}
	return nil
}
