package instances

import (
	"log"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/tkw1536/goprogram/stream"
	"golang.org/x/sync/errgroup"
)

// WissKIInfo represents information about this WissKI Instance.
type WissKIInfo struct {
	Time time.Time // Time this info was built

	// Generic Information
	Slug string // slug
	URL  string // complete URL, including http(s)

	Locked bool // Is this instance currently locked?

	// Information about the running instance
	Running     bool
	LastRebuild time.Time
	LastUpdate  time.Time
	LastCron    time.Time

	// List of backups made
	Snapshots []models.Export

	// WissKI content information
	NoPrefixes   bool              // TODO: Move this into the database
	Prefixes     []string          // list of prefixes
	Pathbuilders map[string]string // all the pathbuilders
}

// Info generate a
func (wisski *WissKI) Info(quick bool) (info WissKIInfo, err error) {
	var group errgroup.Group
	wisski.infoQuick(&info, &group)

	if !quick {
		server, err := wisski.NewPHPServer()
		if err == nil {
			defer server.Close()
		}
		wisski.infoSlow(&info, server, &group)
	}

	err = group.Wait()
	return
}

func (wisski *WissKI) infoQuick(info *WissKIInfo, group *errgroup.Group) {
	info.Time = time.Now().UTC()
	info.Slug = wisski.Slug
	info.URL = wisski.URL().String()

	group.Go(func() (err error) {
		info.Running, err = wisski.Running()
		return
	})

	group.Go(func() (err error) {
		info.Locked = wisski.IsLocked()
		return
	})

	group.Go(func() (err error) {
		info.LastRebuild, _ = wisski.LastRebuild()
		return
	})

	group.Go(func() (err error) {
		info.LastUpdate, _ = wisski.LastUpdate()
		return
	})

	group.Go(func() (err error) {
		info.LastRebuild, _ = wisski.LastRebuild()
		return
	})

	group.Go(func() (err error) {
		info.NoPrefixes = wisski.NoPrefix()
		return
	})
}

func (wisski *WissKI) infoSlow(info *WissKIInfo, server *PHPServer, group *errgroup.Group) {
	group.Go(func() (err error) {
		info.Prefixes, err = wisski.Prefixes(server)
		log.Println("error prefixes: ", err)
		return nil
	})

	group.Go(func() (err error) {
		info.Snapshots, err = wisski.Snapshots()
		log.Println("error snapshots: ", err)
		return nil
	})

	group.Go(func() (err error) {
		info.Pathbuilders, err = wisski.AllPathbuilders(server)
		log.Println("error pathbuilders: ", err)
		return nil
	})

	group.Go(func() (err error) {
		info.LastCron, err = wisski.LastCron(server)
		log.Println("error cron: ", err)
		return
	})
}

// Running checks if this WissKI is currently running.
func (wisski *WissKI) Running() (bool, error) {
	ps, err := wisski.Barrel().Ps(stream.FromNil())
	if err != nil {
		return false, err
	}
	return len(ps) > 0, nil
}
