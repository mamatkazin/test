package users

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"net/http"

// 	"nami/nami_ds/controllers/common"
// )

// func Login(r *http.Request) (data S_Login, err error) {
// 	defer func() {
// 		if rec := recover(); rec != nil {
// 			err = common.GetRecoverError(rec)
// 		}
// 	}()

// 	var login, pass string

// 	if login, pass, err = _checkVariablesLogin(r); err != nil {
// 		panic(err)
// 	}

// 	if data, err = _login(login, pass); err != nil {
// 		if data, err = _login(login, pass); err != nil {
// 			panic(err)
// 		}
// 	}

// 	return
// }

// func _checkVariablesLogin(r *http.Request) (login, pass string, err error) {
// 	defer func() {
// 		if rec := recover(); rec != nil {
// 			err = common.GetRecoverError(rec)
// 		}
// 	}()

// 	login = common.CheckLenString(r.FormValue("login"), 50)
// 	pass = common.CheckLenString(r.FormValue("password"), 50)

// 	return
// }

// func _login(login, pass string) (data S_Login, err error) {
// 	defer func() {
// 		if rec := recover(); rec != nil {
// 			err = common.GetRecoverError(rec)
// 		}
// 	}()

// 	var (
// 		db *sql.DB
// 		tx *sql.Tx

// 		oCenter string
// 	)

// 	if db, err = common.GetDB(0, nil); err != nil {
// 		panic(common.ProcessingError(err.Error()))
// 	}

// 	defer db.Close()

// 	if tx, err = db.Begin(); err != nil {
// 		panic(common.ProcessingError(err.Error()))
// 	}

// 	defer tx.Commit()

// 	if err = tx.QueryRow("select o_User_ID, o_ADT_ID, o_Center from sp_User_Login(?,?)", login, pass).Scan(&data.ID, &data.District.ID, &oCenter); err != nil {
// 		panic(common.ProcessingError(err.Error()))
// 	}

// 	if err = json.Unmarshal([]byte(oCenter), &data.District.LL); err != nil {
// 		panic(common.ProcessingError(err.Error()))
// 	}

// 	return
// }
