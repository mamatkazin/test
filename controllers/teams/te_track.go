package teams

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-pg/pg"

	"nami/nami_ds/controllers/common"
)

func Track(r *http.Request) (data S_Point, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = common.GetRecoverError(rec)
		}
	}()

	var (
		trackID, minute int
	)

	if trackID, minute, err = _checkVariablesTrack(r); err != nil {
		panic(err)
	}

	if data, err = _Track(trackID, minute); err != nil {
		if data, err = _Track(trackID, minute); err != nil {
			panic(err)
		}
	}

	return
}

func _checkVariablesTrack(r *http.Request) (trackID, minute int, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = common.GetRecoverError(rec)
		}
	}()

	// if idStreet, err = strconv.Atoi(r.FormValue("id")); err != nil {
	// 	panic(common.ProcessingError(err.Error()))
	// }

	trackID = 0
	minute = time.Now().Minute()
	//tail = 15

	return
}

func _Track(trackID, minute int) (data S_Point, err error) {
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

	query = "select o_json from nami.fn_track_getpoint(?,?)"
	//query = "select o_json from nami.fn_track_get(?,?,?)"

	if _, err = db.QueryOne(pg.Scan(&oJSON), query, trackID, minute); err != nil {
		panic(common.ProcessingError(err.Error()))
	}

	if err = json.Unmarshal([]byte(oJSON), &data); err != nil {
		panic(common.ProcessingError(err.Error()))
	}

	return
}
