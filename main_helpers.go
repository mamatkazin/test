package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"nami/nami_ds/controllers/common"
)

func initLogFile() (f *os.File, err error) {
	f, err = os.OpenFile("KoT.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		log.Printf("panic: Не возомжно создать лог файл. Ошибка: %s", err.Error())
		return
	}

	log.Println("success: Лог файл успешно создан.")

	// говорим логгеру выводить данные в наш файл
	log.SetOutput(f)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return
}

func initConfigFile() (err error) {
	f, err := os.Open("config.json")

	if err != nil {
		log.Printf("panic: Не возомжно открыть конфигурационный файл. Ошибка: %s", err.Error())
		return
	}

	defer f.Close()

	decoder := json.NewDecoder(f)

	err = decoder.Decode(&common.G_CONFIG)

	if err != nil {
		log.Printf("panic: Не возомжно декодировать конфигурационный файл. Ошибка: %s", err.Error())
		return
	}

	log.Println("success: Конфигурационные данные успешно загружены.")

	return
}

func h_WriteData(access bool, data interface{}, user_id int, w http.ResponseWriter) (err error) {
	var b []byte

	if access {
		if b, err = json.Marshal(data); err != nil {
			panic(common.ProcessingError(err.Error()))
		}

		w.Write(b)
	} else {
		w.WriteHeader(405)
	}

	return
}

func h_Recover(rec interface{}) (res []byte) {
	var (
		err  error
		data = common.S_Data{Valid: false, Errors: make([]string, 0)}
	)

	data.Errors = append(data.Errors, common.GetRecoverErrorText(rec))

	if res, err = json.Marshal(data); err != nil {
		res = []byte(common.ProcessingError(err.Error()))
	}

	return
}
