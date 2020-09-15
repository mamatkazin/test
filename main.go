package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-pg/pg"
	"github.com/gorilla/mux"

	"hospital_track/controllers/common"
)

type S_Status struct {
	Server bool        `json:"server"`
	DB     bool        `json:"db"`
	TI     interface{} `json:"time_idle"`
}

type S_Body struct {
	IMEI      int64   `json:"imei"`
	MAC       string  `json:"mac"`
	TS        int64   `json:"timestamp"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Speed     float64 `json:"speed"`
	Length    float64 `json:"length"`
	Hindrace  int     `json:"hindrace"`
	Direction int     `json:"direction"`
}

type S_Answer_Ins struct {
	Valid bool `json:"valid"`
	//	Hindrace bool  `json:"hindrace"`
	TS    int64 `json:"timestamp"`
	TL    bool  `json:"traffic_light"`
	TLNII bool  `json:"traffic_light_nii"`
}

// type S_Sensor struct {
// 	MAC   string `json:"mac"`
// 	TS    int64  `json:"timestamp"`
// 	Value bool   `json:"value"`
// }

// type S_Answer_Sensor struct {
// 	Valid bool  `json:"valid"`
// 	Value bool  `json:"value"`
// 	TS    int64 `json:"timestamp"`
// }

func init() {
	if err := initConfigFile(); err != nil {
		return
	}
}

// var G_HINDRACE bool

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

	// G_HINDRACE = false

	if err = common.ConnectDB(0); err != nil {
		panic(common.ProcessingError(err.Error()))
	}

	log.Printf("success: port=%s\n", common.G_CONFIG.APP_PORT)

	mx = mux.NewRouter()

	sp := http.StripPrefix("/", http.FileServer(http.Dir("./")))

	mx.HandleFunc("/api/tracks", trackHandler)

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
	var (
		bodyBytes, b []byte
		err          error
		bodySTR      string
		body         S_Body
		oID, oTID    int64
		tm           time.Time
		res          S_Answer_Ins
	)

	defer func() {
		if rec := recover(); rec != nil {
			common.ProcessingError(fmt.Sprintf("Error trackHandler: %v", rec))

			fmt.Println(fmt.Sprintf("Error trackHandler: %v", rec))

			res.Valid = false
			res.TS = time.Now().Unix()

			if b, err = json.Marshal(res); err != nil {
				b = []byte(common.ProcessingError(err.Error()))
			}

			w.Write(b)
		}
	}()

	if r.Method == "POST" {

		if bodyBytes, err = ioutil.ReadAll(r.Body); err != nil {
			panic("Error reading body: " + err.Error())
		}

		bodySTR = string(bodyBytes)

		if err = json.Unmarshal([]byte(bodySTR), &body); err != nil {
			panic(common.ProcessingError(err.Error()))
		}

		//tm = time.Unix(body.TS, 0)
		tm = time.Unix(0, body.TS*int64(time.Millisecond))

		//body = S_Body{4345456544895, "E4:54:E8:1E:62:96", 1578931049289, 37.622504, 55.753215}

		fmt.Println(tm, r.Method, body.MAC, body.TS, body.Lng, body.Lat, body.Speed, body.Length*1000, body.Hindrace, body.Direction)

		if _, err = common.G_INS.QueryOne(pg.Scan(&oID), tm, body.MAC, body.Lng, body.Lat, body.Speed, body.Length*1000, body.Hindrace, body.Direction); err != nil {
			if !common.CheckDB() {
				if err = common.ConnectDB(0); err != nil {
					panic(common.ProcessingError(err.Error()))
				}
			}

			panic(common.ProcessingError(err.Error()))
		}

		if oID > 0 {
			if _, err = common.G_COMPUTED.QueryOne(pg.Scan(&oTID, &res.TL, &res.TLNII), oID); err != nil {
				if !common.CheckDB() {
					if err = common.ConnectDB(0); err != nil {
						panic(common.ProcessingError(err.Error()))
					}
				}
			}

			res.Valid = true

			fmt.Println("oTID", oTID)

			//go func() {
			// if oTID > 0 {
			// 	if _, err = common.G_SPEED.QueryOne(pg.Scan(&oTID2), oTID); err != nil {

			// 		fmt.Println("G_SPEED", err)

			// 		if !common.CheckDB() {
			// 			if err = common.ConnectDB(0); err != nil {
			// 				panic(common.ProcessingError(err.Error()))
			// 			}
			// 		}
			// 	}

			// 	fmt.Println("oTID2", oTID, oTID2)
			// }
			//}()

		} else {
			res.Valid = false
		}

		fmt.Println("Reading oID: ", oID)

		res.TS = time.Now().Unix()

		if b, err = json.Marshal(res); err != nil {
			panic(common.ProcessingError(err.Error()))
		}

		w.Write(b)
	} else {
		res.Valid = false
		res.TS = time.Now().Unix()

		if b, err = json.Marshal(res); err != nil {
			b = []byte(common.ProcessingError(err.Error()))
		}

		w.Write(b)
	}

}
