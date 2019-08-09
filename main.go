package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-pg/pg"

	"github.com/gorilla/mux"

	"nami/nami_track/controllers/common"
)

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

	log.Printf("success: port=%s\n", common.G_CONFIG.APP_PORT)

	mx = mux.NewRouter()

	sp := http.StripPrefix("/", http.FileServer(http.Dir("./")))

	mx.HandleFunc("/api/track", trackHandler)

	mx.PathPrefix("/").Handler(sp)

	if f_log, err = initLogFile(); err != nil {
		return
	}

	defer f_log.Close()

	err = http.ListenAndServe(common.G_CONFIG.APP_PORT, mx)

	if err != nil {
		fmt.Printf("panic: Не возможно запустить сервер. Ошибка: %s", err.Error())
		return
	}
}

func trackHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("Error trackHandler: %v", rec.(error))
		}
	}()

	var (
		bodyBytes []byte
		err       error
		bodySTR   string
		str       []string
		oID       int
	)

	if r.Method == "POST" {

		if bodyBytes, err = ioutil.ReadAll(r.Body); err != nil {
			log.Printf("Error reading body: %v", err)
		}

		bodySTR = string(bodyBytes)

		str = strings.Split(bodySTR, ";")

		//fmt.Println("str", str[1], str[4], str[5], time.Now())

		var (
			db    *pg.DB
			query string
		)

		if db, err = common.GetPGDB(0); err != nil {
			fmt.Println(err.Error())
			panic(common.ProcessingError(err.Error()))
		}

		defer db.Close()

		query = "select o_ID from nami.fn_trackdata_ins(?,?,?,?)"

		if _, err = db.QueryOne(pg.Scan(&oID), query, time.Now(), str[1], str[4], str[5]); err != nil {
			fmt.Println(err.Error())
			panic(common.ProcessingError(err.Error()))
		}

		//fmt.Println("oID", oID)

	}

}
