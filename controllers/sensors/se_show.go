package sensors

import (
	//	"encoding/json"
	//	"io/ioutil"
	"nami/nami_track/controllers/common"
	"net/http"

	"github.com/go-pg/pg"
)

func Show(r *http.Request) (data S_Sensor, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = common.GetRecoverError(rec)
		}
	}()

	var (
		query string
	)

	// if bodyBytes, err = ioutil.ReadAll(r.Body); err != nil {
	// 	panic("Error reading body: " + err.Error())
	// }

	// bodySTR = string(bodyBytes)

	// fmt.Println("@@@@@@@@@", bodySTR)

	// if err = json.Unmarshal([]byte(bodySTR), &body); err != nil {
	// 	panic(common.ProcessingError(err.Error()))
	// }

	query = "select o_UT, o_Flag from nami.fn_hindrace_get(?)"

	// CREATE OR REPLACE FUNCTION nami.fn_hindrace_get (
	//       i_Vid    varchar(30),-- ид устройства
	// out    o_UT     bigint,  -- юникс время прибора,
	// out    o_Flag   boolean  -- признак чистоты зоны
	// )
	// AS
	// $body$

	if _, err = common.G_DB.QueryOne(pg.Scan(&data.Timestamp, &data.Value), query, common.G_CONFIG.SENSOR); err != nil {
		if !common.CheckDB() {
			if err = common.ConnectDB(0); err != nil {
				panic(common.ProcessingError(err.Error()))
			}
		}

		panic(common.ProcessingError(err.Error()))
	}

	data.MAC = common.G_CONFIG.SENSOR

	return

}
