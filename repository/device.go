package repository

import (
	"context"
	"database/sql"
	"math"
	"time"

	"hospital_track/env"
	"hospital_track/models"
)

// SRepositoryDevice структура для работы с базой postgresql
type SRepositoryDevice struct {
	db *sql.DB
}

// PGDevice конструктор
func RepositoryDevice(db *sql.DB) *SRepositoryDevice {
	return &SRepositoryDevice{db: db}
}

// Registry регистрация местоположения устройства в системе
func (db *SRepositoryDevice) Registry(ctx context.Context, device *models.SDevice) (oID int, err error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := "select o_ID from nami.fn_trackdata_ins($1,$2,$3,$4,$5,$6,$7,$8)"

	// CREATE OR REPLACE FUNCTION nami.fn_trackdata_ins (
	// 	i_Time    timestamp,       -- время снятия показаний с прибора
	// 	i_MAC     varchar(30),     -- ид прибора
	// 	i_X       DOUBLE PRECISION,-- долгота
	// 	i_Y       DOUBLE PRECISION,-- широта
	// 	i_Speed   DOUBLE PRECISION,-- скорость
	// 	i_Len     DOUBLE PRECISION,-- длина пути
	// 	i_Dist    INTEGER         ,-- расстояние до помехи
	// 	i_Direct  INTEGER         ,-- направление помехи
	// 	out o_ID  bigint           -- ид точки трека; (-1) если устройства нет в списке, (-2) время бьет назад
	// )

	tm := time.Unix(0, device.TS*int64(time.Millisecond))
	tm = tm.Add(time.Duration(-3 * time.Hour))

	if err = db.db.QueryRowContext(
		ctx,
		query,
		// time.Unix(0, device.TS*int64(time.Millisecond)),
		tm,
		device.MAC,
		device.Lng,
		device.Lat,
		device.Speed,
		device.Length*1000,
		math.Round(device.Hindrace),
		math.Round(device.Direction),
	).Scan(&oID); err != nil {
		if ctx.Err() != nil {
			err = &env.SystemError{Err: ctx.Err().Error()}
		} else {
			err = &env.SystemError{Err: err.Error()}
		}

		env.Logging(err.Error())

		return
	}

	return
}

// Computed посадка точки на граф
func (db *SRepositoryDevice) Computed(ctx context.Context, id int) (oTID int, err error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := "select o_TID from nami.fn_trackdata_computed($1)"

	// CREATE OR REPLACE FUNCTION nami.fn_trackdata_computed (
	// i_TID     bigint,  -- ид точки трека
	// OUT o_TID     bigint   -- ид посаженной точки
	// )

	if err = db.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(&oTID); err != nil {
		if ctx.Err() != nil {
			err = &env.SystemError{Err: ctx.Err().Error()}
		} else {
			err = &env.SystemError{Err: err.Error()}
		}

		env.Logging(err.Error())

		return
	}

	return
}
