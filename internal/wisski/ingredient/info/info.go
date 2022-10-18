package info

import (
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/barrel/drush"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/locker"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php"
	"github.com/FAU-CDI/wisski-distillery/internal/wisski/ingredient/php/extras"
	"golang.org/x/sync/errgroup"
)

type Info struct {
	ingredient.Base

	PHP         *php.PHP
	Barrel      *barrel.Barrel
	Locker      *locker.Locker
	Drush       *drush.Drush
	Prefixes    *extras.Prefixes
	Pathbuilder *extras.Pathbuilder
}

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

// Fetch fetches information about this WissKI.
// TODO: Rework this to be able to determine what kind of information is available.
func (wisski *Info) Fetch(quick bool) (info WissKIInfo, err error) {
	var group errgroup.Group
	wisski.infoQuick(&info, &group)

	if !quick {
		server, err := wisski.PHP.NewServer()
		if err == nil {
			defer server.Close()
		}
		wisski.infoSlow(&info, server, &group)
	}

	err = group.Wait()
	return
}

func (wisski *Info) infoQuick(info *WissKIInfo, group *errgroup.Group) {
	info.Time = time.Now().UTC()
	info.Slug = wisski.Slug
	info.URL = wisski.URL().String()

	group.Go(func() (err error) {
		info.Running, err = wisski.Barrel.Running()
		return
	})

	group.Go(func() (err error) {
		info.Locked = wisski.Locker.Locked()
		return
	})

	group.Go(func() (err error) {
		info.LastRebuild, _ = wisski.Barrel.LastRebuild()
		return
	})

	group.Go(func() (err error) {
		info.LastUpdate, _ = wisski.Drush.LastUpdate()
		return
	})

	group.Go(func() (err error) {
		info.NoPrefixes = wisski.Prefixes.NoPrefix()
		return
	})
}

func (wisski *Info) infoSlow(info *WissKIInfo, server *php.Server, group *errgroup.Group) {
	group.Go(func() (err error) {
		info.Prefixes, _ = wisski.Prefixes.All(server)
		return nil
	})

	group.Go(func() (err error) {
		info.Snapshots, _ = wisski.Snapshots()
		return nil
	})

	group.Go(func() (err error) {
		info.Pathbuilders, _ = wisski.Pathbuilder.GetAll(server)
		return nil
	})

	group.Go(func() (err error) {
		info.LastCron, _ = wisski.Drush.LastCron(server)
		return
	})
}
