package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	//	"strconv"

	"encoding/json"
	//"strings"
	"time"

	"github.com/go-pg/pg"

	"github.com/gorilla/mux"

	"nami/nami_track/controllers/common"
	"nami/nami_track/controllers/sensors"
)

type S_Status struct {
	Server bool        `json:"server"`
	DB     bool        `json:"db"`
	TI     interface{} `json:"time_idle"`
}

type S_Body struct {
	IMEI   int64   `json:"imei"`
	MAC    string  `json:"mac"`
	TS     int64   `json:"timestamp"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Speed  float64 `json:"speed"`
	Length float64 `json:"length"`
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

	//arr := strings.Split(S, "  ")

	if err = common.ConnectDB(0); err != nil {
		panic(common.ProcessingError(err.Error()))
	}

	// for _, item := range arr {
	// 	if _, err = common.G_SPEED.Exec(item); err != nil {
	// 		if !common.CheckDB() {
	// 			if err = common.ConnectDB(0); err != nil {
	// 				panic(common.ProcessingError(err.Error()))
	// 			}
	// 		}
	// 	}
	// }

	// fmt.Println("ok")

	log.Printf("success: port=%s\n", common.G_CONFIG.APP_PORT)

	mx = mux.NewRouter()

	sp := http.StripPrefix("/", http.FileServer(http.Dir("./")))

	mx.HandleFunc("/api/tracks", trackHandler)
	mx.HandleFunc("/api/sensors", sensorsHandler)

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
		body      S_Body
		oID, oTID int64
		tm        time.Time
	)

	//if r.Method == "POST" {

	if bodyBytes, err = ioutil.ReadAll(r.Body); err != nil {
		panic("Error reading body: " + err.Error())
	}

	bodySTR = string(bodyBytes)

	if err = json.Unmarshal([]byte(bodySTR), &body); err != nil {
		panic(common.ProcessingError(err.Error()))
	}

	tm = time.Unix(int64(body.TS/1000), 0)

	//body = S_Body{4345456544895, "E4:54:E8:1E:62:96", 1578931049289, 37.622504, 55.753215}

	fmt.Println("@@@@@@@@@", body.TS, tm, body.Lng, body.Lat, body.Speed, body.Length)

	if _, err = common.G_INS.QueryOne(pg.Scan(&oID), tm, body.MAC, body.Lng, body.Lat, body.Speed, body.Length); err != nil {
		if !common.CheckDB() {
			if err = common.ConnectDB(0); err != nil {
				panic(common.ProcessingError(err.Error()))
			}
		}

		panic(common.ProcessingError(err.Error()))
	}

	if oID > 0 {
		if _, err = common.G_COMPUTED.QueryOne(pg.Scan(&oTID), oID); err != nil {
			if !common.CheckDB() {
				if err = common.ConnectDB(0); err != nil {
					panic(common.ProcessingError(err.Error()))
				}
			}
		}

		//<<<<<<< HEAD
		// go func() {
		// 	if _, err = common.G_SPEED.Exec(oID); err != nil {
		// 		if !common.CheckDB() {
		// 			if err = common.ConnectDB(0); err != nil {
		// 				panic(common.ProcessingError(err.Error()))
		// 			}
		// 		}
		// 	}

		// }()
		// =======
		// 		fmt.Println("oTID", oTID)

		// 		go func() {
		// 			if oTID > 0 {
		// 				if _, err = common.G_SPEED.QueryOne(pg.Scan(&oTID2), oTID); err != nil {

		// 					fmt.Println("G_SPEED", err)

		// 					if !common.CheckDB() {
		// 						if err = common.ConnectDB(0); err != nil {
		// 							panic(common.ProcessingError(err.Error()))
		// 						}
		// 					}
		// 				}

		// 				fmt.Println("oTID2", oTID, oTID2)
		// 			}
		// 		}()
		// >>>>>>> acf5fe408e7015fabf95848eb422a37e59e27a84

		// if insID > 0 {
		// 	go func() {
		// 		if _, err = common.G_SPEED.Exec(insID, raceID); err != nil {
		// 			fmt.Println("G_SPEED", err)

		// 			if !common.CheckDB() {
		// 				if err = common.ConnectDB(0); err != nil {
		// 					common.ProcessingError(err.Error())
		// 				}
		// 			}
		// 		}

		// 	}()
		// }
	}

	fmt.Println("Reading oID: ", oID)

	//w.WriteHeader()

}

func sensorsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	defer func() {
		if rec := recover(); rec != nil {
			w.Write(h_Recover(rec))
		}
	}()

	var (
		err error

		access = true
		//data   = common.S_Data{Valid: true, Errors: make([]string, 0)}
		data interface{}
	)

	if r.Method != "OPTIONS" {
		if r.Method == "GET" {
			if data, err = sensors.Show(r); err != nil {
				panic(err)
			}
		} else if r.Method == "POST" {
			if data, err = sensors.Create(r); err != nil {
				panic(err)
			}

		} else {
			access = false
		}
	}

	if err = h_WriteData(access, data, 0, w); err != nil {
		panic(err)
	}

}
