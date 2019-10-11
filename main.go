package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"strings"
	"time"

	"github.com/go-pg/pg"

	"github.com/gorilla/mux"

	"nami/nami_track/controllers/common"
)

var DB *pg.DB

func init() {
	if err := initConfigFile(); err != nil {
		return
	}
}

func main() {
	var (
		mx    *mux.Router
		err   error
		f_log *os.File
	)

	if DB, err = common.GetPGDB(0); err != nil {
		fmt.Println("GetPGDB", err)
		panic(err.Error())
	}

	defer DB.Close()

	log.Printf("success: port=%s\n", common.G_CONFIG.APP_PORT)

	mx = mux.NewRouter()

	sp := http.StripPrefix("/", http.FileServer(http.Dir("./")))

	mx.HandleFunc("/api/track", trackHandler)
	mx.HandleFunc("/api/track/reserve", trackReserveHandler)

	mx.PathPrefix("/").Handler(sp)

	if f_log, err = initLogFile(); err != nil {
		return
	}

	defer f_log.Close()

	err = http.ListenAndServe(common.G_CONFIG.APP_PORT, mx)

	if err != nil {
		common.ProcessingError("panic: Не возможно запустить сервер. Ошибка: " + err.Error())
		return
	}
}

func trackHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			common.ProcessingError("Error trackHandler: " + rec.(error).Error())
		}
	}()

	var (
		bodyBytes      []byte
		err            error
		bodySTR, query string
		str            []string
		oID, i         int64
		tm             time.Time

		tx *pg.Tx
	)

	//if r.Method == "POST" {

	if bodyBytes, err = ioutil.ReadAll(r.Body); err != nil {
		panic("Error reading body: " + err.Error())
	}

	bodySTR = string(bodyBytes)

	str = strings.Split(bodySTR, ";")

	if i, err = strconv.ParseInt(str[0], 10, 64); err != nil {
		panic(err.Error())
	}

	tm = time.Unix(i, 0)

	if tx, err = DB.Begin(); err != nil {
		panic(err.Error())
	}

	query = "select o_ID from nami.fn_trackdata_ins(?,?,?,?)"

	// CREATE OR REPLACE FUNCTION nami.fn_trackdata_ins (
	//   i_Time    timestamp,       -- время снятия показаний с прибора
	//   i_MAC     varchar(30),     -- мак прибора
	//   i_X       DOUBLE PRECISION,-- долгота
	//   i_Y       DOUBLE PRECISION,-- широта
	//   out o_ID  bigint           -- ид точки трека; (-1) если устройства нет в списке, (-2) время бьет назад
	// )
	// AS

	fmt.Println(tm, str[1], str[5], str[4])

	if _, err = tx.QueryOne(pg.Scan(&oID), query, tm, str[1], str[5], str[4]); err != nil {
		fmt.Println("QueryOne", err)
		tx.Rollback()
		panic(err.Error())
	}

	tx.Commit()

	if oID > 0 {
		if tx, err = DB.Begin(); err != nil {
			panic(err.Error())
		}

		query = "select from nami.fn_trackdata_computed(?)"

		// CREATE OR REPLACE FUNCTION nami.fn_trackdata_computed (
		//   i_TID     bigint  -- ид точки трека
		// )

		if _, err = tx.Exec(query, oID); err != nil {
			fmt.Println("Exec", err)
			tx.Rollback()
			panic(err.Error())
		}

		tx.Commit()
	}

	fmt.Println("Reading oID: ", oID)

	//}

}

func trackReserveHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			common.ProcessingError("Error trackHandler: " + rec.(error).Error())
		}
	}()

	// var (
	// 	bodyBytes      []byte
	// 	err            error
	// 	bodySTR, query string
	// 	str            []string
	// 	oID, i         int64
	// 	tm             time.Time

	// 	db *pg.DB
	// 	tx *pg.Tx
	// )

	if r.Method == "GET" {

		fmt.Println("MNN", r.FormValue("date_time"), r.FormValue("lat"), r.FormValue("lng"))

		// if db, err = common.GetPGDB(0); err != nil {
		// 	fmt.Println("GetPGDB", err)
		// 	panic(err.Error())
		// }

		// defer db.Close()

		// if tx, err = db.Begin(); err != nil {
		// 	fmt.Println("Begin", err)
		// 	panic(err.Error())
		// }

		// defer tx.Commit()

		// query = "select o_ID from nami.fn_trackdata_ins(?,?,?,?)"

		// // CREATE OR REPLACE FUNCTION nami.fn_trackdata_ins (
		// //   i_Time    timestamp,       -- время снятия показаний с прибора
		// //   i_MAC     varchar(30),     -- мак прибора
		// //   i_X       DOUBLE PRECISION,-- долгота
		// //   i_Y       DOUBLE PRECISION,-- широта
		// //   out o_ID  bigint           -- ид точки трека; (-1) если устройства нет в списке, (-2) время бьет назад
		// // )
		// // AS

		// fmt.Println(tm, str[1], str[5], str[4])

		// if _, err = db.QueryOne(pg.Scan(&oID), query, tm, str[1], str[5], str[4]); err != nil {
		// 	fmt.Println("QueryOne", err)
		// 	tx.Rollback()
		// 	panic(err.Error())
		// }

		// log.Printf("Reading oID: %v", oID)

	}

}
