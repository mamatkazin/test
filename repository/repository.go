package repository

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"hospital_track/models"
)

// IDevice интерфейс с методами устройства
type IDevice interface {
	Registry(ctx context.Context, device *models.SDevice) (int, error)
	Computed(ctx context.Context, id int) (int, error)
}

// SRepository структура репозитория
type SRepository struct {
	IDevice
}

// Repository конструктор репозитория
func Repository(strConn string) (rep *SRepository, err error) {
	var db *sql.DB

	if db, err = pgOpen(strConn); err != nil {
		return
	}

	rep = &SRepository{
		IDevice: RepositoryDevice(db),
	}

	return
}

func pgOpen(strConn string) (db *sql.DB, err error) {
	if db, err = sql.Open("postgres", strConn); err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return
	}

	return

}
