package exporter

import "github.com/FAU-CDI/wisski-distillery/internal/models"

func (backup *Backup) LogEntry() models.Export {
	return models.Export{
		Created: backup.StartTime,
		Slug:    "",
	}
}

func (snapshot *Snapshot) LogEntry() models.Export {
	return models.Export{
		Created: snapshot.StartTime,
		Slug:    snapshot.Instance.Slug,
	}
}
