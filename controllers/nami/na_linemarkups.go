package nami

import (
	"encoding/json"
	"net/http"

	"github.com/go-pg/pg"

	"nami/nami_ds/controllers/common"
)

func Linemarkups(r *http.Request) (data S_MultiPolygon, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = common.GetRecoverError(rec)
		}
	}()

	if data, err = _Linemarkups(); err != nil {
		if data, err = _Linemarkups(); err != nil {
			panic(err)
		}
	}

	return
}

func _Linemarkups() (data S_MultiPolygon, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = common.GetRecoverError(rec)
		}
	}()

	var (
		db           *pg.DB
		oJSON, query string
	)

	if db, err = common.GetPGDB(0); err != nil {
		panic(common.ProcessingError(err.Error()))
	}

	defer db.Close()

	query = "select o_FC_Json from nami.fn_roadmarkup_getcollection()"

	if _, err = db.QueryOne(pg.Scan(&oJSON), query); err != nil {
		panic(common.ProcessingError(err.Error()))
	}

	if err = json.Unmarshal([]byte(oJSON), &data); err != nil {
		panic(common.ProcessingError(err.Error()))
	}

	return
}
