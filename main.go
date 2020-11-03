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

	log.Printf("success: port=%s\n", env.GCONFIG.AppPort)

	if err = http.Run(env.GCONFIG.AppPort); err != nil {
		log.Printf("panic: Не возомжно стартануть сервер. Ошибка: %s", err.Error())
		return
	}
}

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
