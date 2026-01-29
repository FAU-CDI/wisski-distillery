package exporter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/FAU-CDI/wisski-distillery/internal/models"
)

// snapshotJSON is the JSON representation of [Snapshot].
type snapshotJSON struct {
	Description SnapshotDescription `json:"Description"`
	Instance    models.Instance     `json:"Instance"`

	StartTime time.Time `json:"StartTime"`
	EndTime   time.Time `json:"EndTime"`

	ErrPanic json.RawMessage            `json:"ErrPanic,omitempty"`
	ErrStart json.RawMessage            `json:"ErrStart,omitempty"`
	ErrStop  json.RawMessage            `json:"ErrStop,omitempty"`
	Errors   map[string]json.RawMessage `json:"Errors,omitempty"`

	Logs     map[string]string `json:"Logs,omitempty"`
	Manifest []string          `json:"Manifest,omitempty"`
}

func (s Snapshot) MarshalJSON() ([]byte, error) {
	j := snapshotJSON{
		Description: s.Description,
		Instance:    s.Instance,
		StartTime:   s.StartTime,
		EndTime:     s.EndTime,
		Logs:        s.Logs,
		Manifest:    s.Manifest,
	}

	// marshal all the error fields as json strings.
	if s.ErrPanic != nil {
		j.ErrPanic = marshalString(fmt.Sprint(s.ErrPanic))
	}
	if s.ErrStart != nil {
		j.ErrStart = marshalString(s.ErrStart.Error())
	}
	if s.ErrStop != nil {
		j.ErrStop = marshalString(s.ErrStop.Error())
	}
	if len(s.Errors) > 0 {
		j.Errors = make(map[string]json.RawMessage, len(s.Errors))
		for k, v := range s.Errors {
			if v != nil {
				j.Errors[k] = marshalString(v.Error())
			}
		}
	}

	return json.Marshal(j)
}

func marshalString(s string) json.RawMessage {
	// we are only marshalling a string.
	// which cannot cause any error.
	value, _ := json.Marshal(s)
	return value
}

func (s *Snapshot) UnmarshalJSON(data []byte) error {
	var j snapshotJSON
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}

	// unmarshal all the fields, but skip the error fields.
	// this prevents errors while unmarshalling interfaces.

	s.Description = j.Description
	s.Instance = j.Instance
	s.StartTime = j.StartTime
	s.EndTime = j.EndTime
	s.Logs = j.Logs
	s.Manifest = j.Manifest

	return nil
}
