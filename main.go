package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"nami/nami_ds/controllers/common"
	"nami/nami_ds/controllers/cookies"
	"nami/nami_ds/controllers/nami"
	"nami/nami_ds/controllers/teams"
	//	"nami/nami_ds/controllers/users"
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

	mx.HandleFunc("/", indexHandler)
	mx.HandleFunc("/logout", logoutHandler)

	mx.HandleFunc("/api/roads", roadsHandler)
	mx.HandleFunc("/api/linemarkups", linemarkupsHandler)
	mx.HandleFunc("/api/bufferzone", bufferzoneHandler)
	mx.HandleFunc("/api/teams/{t_id}/track", teamsTrackHandler)

	mx.HandleFunc("/nami", indexHandler)
	mx.HandleFunc("/nami/", indexHandler)
	mx.HandleFunc("/nami/{page}", indexHandler)
	mx.HandleFunc("/nami/{page}/", indexHandler)
	mx.HandleFunc("/nami/{page}/{page}", indexHandler)
	mx.HandleFunc("/nami/{page}/{page}/", indexHandler)
	mx.HandleFunc("/nami/{page}/{page}/{page}", indexHandler)
	mx.HandleFunc("/nami/{page}/{page}/{page}/", indexHandler)
	mx.HandleFunc("/nami/{page}/{page}/{page}/{page}", indexHandler)
	mx.HandleFunc("/nami/{page}/{page}/{page}/{page}/", indexHandler)
	mx.HandleFunc("/nami/{page}/{page}/{page}/{page}/{page}", indexHandler)
	mx.HandleFunc("/nami/{page}/{page}/{page}/{page}/{page}/", indexHandler)
	mx.HandleFunc("/nami/{page}/{page}/{page}/{page}/{page}/{page}", indexHandler)
	mx.HandleFunc("/nami/{page}/{page}/{page}/{page}/{page}/{page}/", indexHandler)
	mx.HandleFunc("/nami/{page}/{page}/{page}/{page}/{page}/{page}/{page}", indexHandler)
	mx.HandleFunc("/nami/{page}/{page}/{page}/{page}/{page}/{page}/{page}/", indexHandler)

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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	defer func() {
		if rec := recover(); rec != nil {
			w.Write(h_Recover(rec))
		}
	}()

	var (
		page *template.Template
		err  error
	)

	if page, err = template.ParseFiles("app/index.html"); err != nil {
		panic(common.ProcessingError(err.Error()))
	}

	page.ExecuteTemplate(w, "index", nil)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookieData := new(common.CookieData)
	cookieData.UserID = -1
	cookies.SetSession(*cookieData, w)
	cookies.ClearSession(w)
	http.Redirect(w, r, "/", 302)
}

func roadsHandler(w http.ResponseWriter, r *http.Request) {
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

	cookieData = cookies.GetSession(r)

	if common.G_CONFIG.RUN_MODE == "DEV" {
		cookieData.UserID = 0
	}

	if cookieData.UserID != -1 {
		if r.Method != "OPTIONS" {
			if r.Method == "GET" {
				if data.Items, err = nami.Roads(r); err != nil {
					panic(err)
				}
			} else {
				access = false
			}
		}

		if err = h_WriteData(access, data, cookieData.UserID, w); err != nil {
			panic(err)
		}

	} else {
		w.WriteHeader(401)
	}
}

func linemarkupsHandler(w http.ResponseWriter, r *http.Request) {
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

	cookieData = cookies.GetSession(r)

	if common.G_CONFIG.RUN_MODE == "DEV" {
		cookieData.UserID = 0
	}

	if cookieData.UserID != -1 {
		if r.Method != "OPTIONS" {
			if r.Method == "GET" {
				if data.Items, err = nami.Linemarkups(r); err != nil {
					panic(err)
				}
			} else {
				access = false
			}
		}

		if err = h_WriteData(access, data, cookieData.UserID, w); err != nil {
			panic(err)
		}

	} else {
		w.WriteHeader(401)
	}
}

func bufferzoneHandler(w http.ResponseWriter, r *http.Request) {
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

	cookieData = cookies.GetSession(r)

	if common.G_CONFIG.RUN_MODE == "DEV" {
		cookieData.UserID = 0
	}

	if cookieData.UserID != -1 {
		if r.Method != "OPTIONS" {
			if r.Method == "GET" {
				if data.Items, err = nami.Buffers(r); err != nil {
					panic(err)
				}
			} else {
				access = false
			}
		}

		if err = h_WriteData(access, data, cookieData.UserID, w); err != nil {
			panic(err)
		}

	} else {
		w.WriteHeader(401)
	}
}

func teamsTrackHandler(w http.ResponseWriter, r *http.Request) {
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

	cookieData = cookies.GetSession(r)

	if common.G_CONFIG.RUN_MODE == "DEV" {
		cookieData.UserID = 0
	}

	if cookieData.UserID != -1 {
		if r.Method != "OPTIONS" {
			if r.Method == "GET" {
				if data.Items, err = teams.Track(r); err != nil {
					panic(err)
				}
			} else {
				access = false
			}
		}

		if err = h_WriteData(access, data, cookieData.UserID, w); err != nil {
			panic(err)
		}

	} else {
		w.WriteHeader(401)
	}
}
