package snapshots

import "github.com/FAU-CDI/wisski-distillery/internal/models"

func (backup *Backup) LogEntry() models.Snapshot {
	return models.Snapshot{
		Created: backup.StartTime,
		Slug:    "",
	}
}

func (snapshot *Snapshot) LogEntry() models.Snapshot {
	return models.Snapshot{
		Created: snapshot.StartTime,
		Slug:    snapshot.Instance.Slug,
	}
}
