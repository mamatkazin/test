package sensors

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"nami/nami_track/controllers/common"
	"net/http"

	"github.com/go-pg/pg"
)

func Create(r *http.Request) (id common.S_ID, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = common.GetRecoverError(rec)
		}
	}()

	var (
		bodyBytes      []byte
		bodySTR, query string
		body           S_Sensor
	)

	if bodyBytes, err = ioutil.ReadAll(r.Body); err != nil {
		panic("Error reading body: " + err.Error())
	}

	bodySTR = string(bodyBytes)

	fmt.Println("@@@@@@@@@", bodySTR)

	if err = json.Unmarshal([]byte(bodySTR), &body); err != nil {
		panic(common.ProcessingError(err.Error()))
	}

	query = "select o_ID from nami.fn_hindrace_ins(?,?,?)"

	// CREATE OR REPLACE FUNCTION nami.fn_hindrace_ins (
	//    i_UT     bigint,     -- юникс время прибора,
	//    i_Flag   boolean,    -- признак чистоты зоны,
	//    i_Vid    varchar(30),-- ид устройства
	// out o_ID     bigint
	// )
	// AS

	if _, err = common.G_DB.QueryOne(pg.Scan(&id.ID), query, body.Timestamp, body.Value, body.MAC); err != nil {
		if !common.CheckDB() {
			if err = common.ConnectDB(0); err != nil {
				panic(common.ProcessingError(err.Error()))
			}
		}

		panic(common.ProcessingError(err.Error()))
	}

	fmt.Println("Reading oID: ", id.ID)

	return

}
