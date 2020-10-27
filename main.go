package main

import (
	"encoding/json"
	"log"
	"os"

	"hospital_track/env"
	"hospital_track/handler"
	"hospital_track/repository"
	"hospital_track/service"
)

// type S_Status struct {
// 	Server bool        `json:"server"`
// 	DB     bool        `json:"db"`
// 	TI     interface{} `json:"time_idle"`
// }

// type S_Body struct {
// 	IMEI      int64   `json:"imei"`
// 	MAC       string  `json:"mac"`
// 	TS        int64   `json:"timestamp"`
// 	Lat       float64 `json:"lat"`
// 	Lng       float64 `json:"lng"`
// 	Speed     float64 `json:"speed"`
// 	Length    float64 `json:"length"`
// 	Hindrace  float64 `json:"distObject"`
// 	Direction float64 `json:"ori"`
// }

// type S_Answer_Ins struct {
// 	Valid bool `json:"valid"`
// 	//	Hindrace bool  `json:"hindrace"`
// 	TS    int64 `json:"timestamp"`
// 	TL    bool  `json:"traffic_light"`
// 	TLNII bool  `json:"traffic_light_nii"`
// }

func init() {
	if err := initConfigFile(); err != nil {
		return
	}
}

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("panic: Ошибка: %v", rec)
		}
	}()

	var (
		//mx    *mux.Router
		err   error
		repos *repository.SRepository
	)

	if repos, err = repository.Repository(env.GCONFIG.DBConnect); err != nil {
		log.Printf("panic: Не возомжно подключится к серверу баз данных. Ошибка: %s", err.Error())
		return
	}

	services := service.Service(repos)
	handlers := handler.Handler(services)

	http := handlers.InitRoutes()

	// if err = common.ConnectDB(0); err != nil {
	// 	panic(err.Error())
	// }

	log.Printf("success: port=%s\n", env.GCONFIG.AppPort)

	// mx = mux.NewRouter()

	// sp := http.StripPrefix("/", http.FileServer(http.Dir("./")))

	// mx.HandleFunc("/api/tracks", trackHandler)

	// mx.PathPrefix("/").Handler(sp)

	if err = http.Run(env.GCONFIG.AppPort); err != nil {
		log.Printf("panic: Не возомжно стартануть сервер. Ошибка: %s", err.Error())
		return
	}
}

// func trackHandler(w http.ResponseWriter, r *http.Request) {
// 	var (
// 		bodyBytes []byte
// 		err       error
// 		bodySTR   string
// 		body      S_Body
// 		oID, oTID int64
// 		tm        time.Time
// 		//res          S_Answer_Ins
// 	)

// 	defer func() {
// 		if rec := recover(); rec != nil {

// 			// fmt.Println(fmt.Sprintf("Error trackHandler: %v", rec))

// 			// res.Valid = false
// 			// res.TS = time.Now().Unix()

// 			// if b, err = json.Marshal(res); err != nil {
// 			// 	b = []byte(common.ProcessingError(err.Error()))
// 			// }

// 			// w.Write(b)
// 		}
// 	}()

// 	if r.Method == "POST" {

// 		if bodyBytes, err = ioutil.ReadAll(r.Body); err != nil {
// 			panic("Error reading body: " + err.Error())
// 		}

// 		bodySTR = string(bodyBytes)

// 		if err = json.Unmarshal([]byte(bodySTR), &body); err != nil {
// 			panic(err.Error())
// 		}

// 		tm = time.Unix(0, body.TS*int64(time.Millisecond))

// 		fmt.Println(tm, r.Method, body.MAC, body.TS, body.Lng, body.Lat, body.Speed, body.Length*1000, math.Round(body.Hindrace), math.Round(body.Direction))

// 		if _, err = common.G_INS.QueryOne(pg.Scan(&oID), tm, body.MAC, body.Lng, body.Lat, body.Speed, body.Length*1000, math.Round(body.Hindrace), math.Round(body.Direction)); err != nil {
// 			if !common.CheckDB() {
// 				if err = common.ConnectDB(0); err != nil {
// 					panic(err.Error())
// 				}
// 			}

// 			panic(err.Error())
// 		}

// 		if oID > 0 {
// 			if _, err = common.G_COMPUTED.QueryOne(pg.Scan(&oTID), oID); err != nil {
// 				if !common.CheckDB() {
// 					if err = common.ConnectDB(0); err != nil {
// 						panic(err.Error())
// 					}
// 				}

// 				panic(err.Error())
// 			}

// 			// res.Valid = true

// 			fmt.Println("oTID", oTID)

// 		} else {
// 			//res.Valid = false
// 		}

// 		fmt.Println("Reading oID: ", oID)

// 		// res.TS = time.Now().Unix()

// 		// if b, err = json.Marshal(res); err != nil {
// 		// 	panic(common.ProcessingError(err.Error()))
// 		// }

// 		// w.Write(b)
// 	} else {
// 		// res.Valid = false
// 		// res.TS = time.Now().Unix()

// 		// if b, err = json.Marshal(res); err != nil {
// 		// 	b = []byte(common.ProcessingError(err.Error()))
// 		// }

// 		// w.Write(b)
// 	}

// }

func initConfigFile() (err error) {
	f, err := os.Open("config.json")

	if err != nil {
		log.Printf("panic: Не возомжно открыть конфигурационный файл. Ошибка: %s", err.Error())
		return
	}

	defer f.Close()

	decoder := json.NewDecoder(f)

	err = decoder.Decode(&env.GCONFIG)

	if err != nil {
		log.Printf("panic: Не возомжно декодировать конфигурационный файл. Ошибка: %s", err.Error())
		return
	}

	log.Println("success: Конфигурационные данные успешно загружены.")

	return
}
