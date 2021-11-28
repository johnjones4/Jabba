package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/johnjones4/Jabba/core"
)

type Scannable interface {
	Scan(dest ...interface{}) (err error)
}

type PGStore struct {
	pool       *pgxpool.Pool
	logDirPath string
}

func NewPGStore(url string, logDirPath string) (*PGStore, error) {
	pool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return &PGStore{pool, logDirPath}, nil
}

func (s *PGStore) SaveEvent(event *core.Event) error {
	vendorInfoJson, err := json.Marshal(event.VendorInfo)
	if err != nil {
		return err
	}

	alertsJson, err := json.Marshal(event.Alerts)
	if err != nil {
		return err
	}

	err = s.pool.QueryRow(context.Background(), "INSERT INTO events (event_vendor_type, event_vendor_id, created, vendor_info, alerts, is_normal) VALUES ($1,$2,$3,$4,$5,$6) RETURNING \"id\"",
		event.EventVendorType,
		event.EventVendorID,
		event.Created,
		vendorInfoJson,
		alertsJson,
		event.IsNormal,
	).Scan(&event.ID)
	if err != nil {
		return err
	}

	if event.Log != "" {
		logPath := path.Join(s.logDirPath, fmt.Sprintf("%d.log", event.ID))
		err = os.WriteFile(logPath, []byte(event.Log), 0777)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *PGStore) GetEvents(limit int, offset int) ([]core.Event, error) {
	rows, err := s.pool.Query(context.Background(), "SELECT id, event_vendor_type, event_vendor_id, created, vendor_info, alerts, is_normal FROM events ORDER BY created DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	events := make([]core.Event, 0)
	for rows.Next() {
		event, err := parseEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *PGStore) GetEvent(id int) (core.Event, error) {
	row := s.pool.QueryRow(context.Background(), "SELECT id, event_vendor_type, event_vendor_id, created, vendor_info, alerts, is_normal FROM events WHERE id = $1",
		id,
	)

	event, err := parseEvent(row)
	if err != nil {
		return core.Event{}, err
	}

	logPath := path.Join(s.logDirPath, fmt.Sprintf("%d.log", event.ID))
	_, err = os.Stat(logPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log, err := os.ReadFile(logPath)
		if err != nil {
			return core.Event{}, err
		}

		event.Log = string(log)
	} else if err != nil {
		return core.Event{}, nil
	}

	return event, nil
}

func parseEvent(row Scannable) (core.Event, error) {
	e := core.Event{}
	var vendorInfoBytes []byte
	var alertsBytes []byte
	err := row.Scan(
		&e.ID,
		&e.EventVendorType,
		&e.EventVendorID,
		&e.Created,
		&vendorInfoBytes,
		&alertsBytes,
		&e.IsNormal,
	)
	if err != nil {
		return core.Event{}, err
	}

	err = json.Unmarshal(vendorInfoBytes, &e.VendorInfo)
	if err != nil {
		return core.Event{}, err
	}

	err = json.Unmarshal(alertsBytes, &e.Alerts)
	if err != nil {
		return core.Event{}, err
	}

	return e, nil
}
