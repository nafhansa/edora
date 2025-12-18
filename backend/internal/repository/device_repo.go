package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"edora/backend/internal/models"
)

type DeviceRepository struct {
	db interface{}
}

func NewDeviceRepository(db interface{}) *DeviceRepository {
	return &DeviceRepository{db: db}
}

// DeviceRepo defines methods used by services and tests.
type DeviceRepo interface {
	GetBySerial(ctx context.Context, serial string) (*models.Device, error)
	UpdateLastSeen(ctx context.Context, id string, t time.Time) error
	CountActive(ctx context.Context, since time.Duration) (int, error)
}

func (d *DeviceRepository) GetBySerial(ctx context.Context, serial string) (*models.Device, error) {
	if d.db == nil {
		return nil, nil
	}
	db, ok := d.db.(*sql.DB)
	if !ok {
		return nil, errors.New("unsupported db type")
	}
	q := `SELECT id, serial_number, name, status, last_seen, created_at FROM devices WHERE serial_number = $1 LIMIT 1`
	row := db.QueryRowContext(ctx, q, serial)
	var dev models.Device
	var lastSeen sql.NullTime
	if err := row.Scan(&dev.ID, &dev.SerialNumber, &dev.Name, &dev.Status, &lastSeen, &dev.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if lastSeen.Valid {
		dev.LastSeen = lastSeen.Time
	}
	return &dev, nil
}

func (d *DeviceRepository) UpdateLastSeen(ctx context.Context, id string, t time.Time) error {
	if d.db == nil {
		return nil
	}
	db, ok := d.db.(*sql.DB)
	if !ok {
		return errors.New("unsupported db type")
	}
	q := `UPDATE devices SET last_seen = $1, status = 'online' WHERE id = $2`
	_, err := db.ExecContext(ctx, q, t, id)
	return err
}

// CountActive returns number of devices with last_seen >= threshold duration ago
func (d *DeviceRepository) CountActive(ctx context.Context, since time.Duration) (int, error) {
	if d.db == nil {
		return 0, nil
	}
	db, ok := d.db.(*sql.DB)
	if !ok {
		return 0, errors.New("unsupported db type")
	}
	var cnt int
	q := `SELECT COUNT(*) FROM devices WHERE last_seen >= now() - $1::interval`
	// convert duration to seconds string
	interval := fmt.Sprintf("%d seconds", int64(since.Seconds()))
	if err := db.QueryRowContext(ctx, q, interval).Scan(&cnt); err != nil {
		return 0, err
	}
	return cnt, nil
}
