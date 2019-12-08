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

type S_Status struct {
	Server bool        `json:"server"`
	DB     bool        `json:"db"`
	TI     interface{} `json:"time_idle"`
}

func init() {
	if err := initConfigFile(); err != nil {
		return
	}
}

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			common.ProcessingError("Error main: " + common.GetRecoverErrorText(rec))
		}
	}()

	var (
		mx    *mux.Router
		err   error
		f_log *os.File
	)

	if err = common.ConnectDB(0); err != nil {
		panic(common.ProcessingError(err.Error()))
	}

	log.Printf("success: port=%s\n", common.G_CONFIG.APP_PORT)

	mx = mux.NewRouter()

	sp := http.StripPrefix("/", http.FileServer(http.Dir("./")))

	mx.HandleFunc("/api/track", trackHandler)
	mx.HandleFunc("/api/servers/ping", serversPingHandler)

	mx.PathPrefix("/").Handler(sp)

	if f_log, err = initLogFile(); err != nil {
		panic(err.Error())
	}

	defer f_log.Close()

	err = http.ListenAndServe(common.G_CONFIG.APP_PORT, mx)

	if err != nil {
		panic(err.Error())
	}
}

func trackHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			common.ProcessingError("Error trackHandler: " + rec.(error).Error())
		}
	}()

	var (
		bodyBytes []byte
		err       error
		bodySTR   string
		str       []string
		oID, i    int64
		tm        time.Time
	)

	if r.Method == "POST" {

		if bodyBytes, err = ioutil.ReadAll(r.Body); err != nil {
			panic("Error reading body: " + err.Error())
		}

		bodySTR = string(bodyBytes)

		str = strings.Split(bodySTR, ";")

		if i, err = strconv.ParseInt(str[0], 10, 64); err != nil {
			panic(err.Error())
		}

		tm = time.Unix(i, 0)

		fmt.Println(tm, str[1], str[5], str[4])

		if _, err = common.G_INS.QueryOne(pg.Scan(&oID), tm, str[1], str[5], str[4]); err != nil {
			fmt.Println("QueryOne", err)

			if !common.CheckDB() {
				if err = common.ConnectDB(0); err != nil {
					panic(common.ProcessingError(err.Error()))
				}
			}
		}

		if oID > 0 {
			if _, err = common.G_COMPUTED.Exec(oID); err != nil {
				fmt.Println("Exec", err)

				if !common.CheckDB() {
					if err = common.ConnectDB(0); err != nil {
						panic(common.ProcessingError(err.Error()))
					}
				}
			}

			go func() {
				if _, err = common.G_SPEED.Exec(oID); err != nil {
					fmt.Println("G_SPEED", err)

					if !common.CheckDB() {
						if err = common.ConnectDB(0); err != nil {
							common.ProcessingError(err.Error())
						}
					}
				}

			}()
		}

		fmt.Println("Reading oID: ", oID)
	}

}

func serversPingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	defer func() {
		if rec := recover(); rec != nil {
			w.Write(h_Recover(rec))
		}
	}()

	var (
		err        error
		cookieData common.CookieData

		access = true
		data   = common.S_Data{Valid: true, Errors: make([]string, 0)}
	)

	if r.Method != "OPTIONS" {
		if r.Method == "GET" {
			if data.Items, err = Ping(); err != nil {
				panic(err)
			}

		} else {
			access = false
		}
	}

	if err = h_WriteData(access, data, cookieData.UserID, w); err != nil {
		panic(err)
	}
}

func Ping() (data S_Status, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = common.GetRecoverError(rec)
		}
	}()

	var n int

	data.Server = true
	data.DB = true

	if _, err = common.G_DB.QueryOne(pg.Scan(&n), "SELECT 110+1"); err != nil {
		data.DB = false
		err = nil
	}

	return
}
