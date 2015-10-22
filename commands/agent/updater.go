package agent

import (
	"log"
	"sync"
	"time"

	"github.com/ChrisMcKenzie/dropship/repo"
	"github.com/ChrisMcKenzie/dropship/structs"
)

type updater struct {
	ticker     *time.Ticker
	shutdownCh <-chan struct{}
	service    structs.Service
	repo       repo.Repo
}

func (u *updater) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	log.Println("Starting", u.service.Id, "updater")
	for {
		select {
		case <-u.ticker.C:
			log.Println("Performing", u.service.Id, "update check")
			u.check()
		case _, ok := <-u.shutdownCh:
			if !ok {
				log.Println("Stopping", u.service.Id, "update check")
				u.ticker.Stop()
				return
			}
		}
	}
}

func (u *updater) check() {
	isUpToDate, err := u.repo.IsUpdated(u.service)

	if err != nil {
		log.Fatal(err)
	}

	// check the md5sums
	if !isUpToDate {
		u.update()
	}

	return
}

func (u *updater) update() {
	log.Println("Starting update")
	lock, err := AcquireLock(u.service)
	_, err = lock.Lock(nil)
	if err != nil {
		panic(err)
	}

	file, meta, err := u.repo.Download(u.service)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Finished Downloading")
	if meta.ContentType == "application/x-gzip" {
		err := untar(file, u.service.Artifact.Dest)
		if err != nil {
			log.Fatal(err)
		}
	}

	u.service.Hash = meta.Hash
	log.Println("Setting current version to", u.service.Hash)
	lock.Unlock()
}