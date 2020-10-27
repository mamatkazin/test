package common

import (
	"fmt"
	"hospital_track/env"

	"github.com/go-pg/pg"
)

var (
	G_DB              *pg.DB
	G_INS, G_COMPUTED *pg.Stmt
)

func ConnectDB(count int) (err error) {
	defer func() {
		if rec := recover(); rec != nil {

		}
	}()

	var (
		opt pg.Options
		n   int
	)

	opt.Addr = env.GCONFIG.PG_ADDRESS
	opt.Database = env.GCONFIG.PG_BASE
	opt.User = env.GCONFIG.PG_USER
	opt.Password = env.GCONFIG.PG_PASSWORD

	G_DB = nil

	if count > 10 {
		panic("Потеряно соединение с БД")
	}

	G_DB = pg.Connect(&opt)

	if _, err = G_DB.QueryOne(pg.Scan(&n), "SELECT 110+1"); err != nil {
		fmt.Println(err)
		count = count + 1
		err = ConnectDB(count)

		return
	}

	if err = PrepareStmt(); err != nil {
		fmt.Println(err)
		count = count + 1
		err = ConnectDB(count)

		return
	}

	return
}

func PrepareStmt() (err error) {
	defer func() {
		if rec := recover(); rec != nil {

		}
	}()

	if G_INS, err = G_DB.Prepare("select o_ID from nami.fn_trackdata_ins($1,$2,$3,$4,$5,$6,$7,$8)"); err != nil {
		panic(err.Error())
	}

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

	if G_COMPUTED, err = G_DB.Prepare("select o_TID from nami.fn_trackdata_computed($1)"); err != nil {
		panic(err.Error())
	}

	// 	CREATE OR REPLACE FUNCTION nami.fn_trackdata_computed (
	// 		i_TID     bigint,  -- ид точки трека
	// OUT  o_TID     bigint   -- ид посаженной точки
	// )
	// AS

	return
}

func CheckDB() bool {
	var (
		err error
		n   int
	)

	if G_DB != nil {
		if _, err = G_DB.QueryOne(pg.Scan(&n), "SELECT 110+1"); err == nil {
			return true
		}
	}

	return false
}
