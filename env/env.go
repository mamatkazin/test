package env

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

type Configuration struct {
	AppPort     string `json:"appPort"`
	DBConnect   string `json:"dbConnect"`
	MailEmail   string `json:"mailEmail"`
	MailPass    string `json:"mailPass"`
	MailSmtp    string `json:"mailSmtp"`
	MailPort    string `json:"mailPort"`
	MailDevs    string `json:"mailDevs"`
	PG_ADDRESS  string
	PG_USER     string
	PG_PASSWORD string
	PG_BASE     string
}

//OutputData структура для вывода данных
type OutputData struct {
	Valid  bool        `json:"valid"`
	Errors string      `json:"errors"`
	Items  interface{} `json:"items"`
}

type Env struct {
	DataBases map[string]*sql.DB
}

func NewEnv() *Env {
	return &Env{
		DataBases: make(map[string]*sql.DB),
	}
}

var (
	GCONFIG Configuration
)

func SendDeveloper(message_type, text string) {
	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: error = %v", rec)
			}
		}()

		var (
			arr_to []string
			m      *gomail.Message
			d      *gomail.Dialer
			err    error
			port   int
		)

		if GCONFIG.MailDevs != "" {
			arr_to = strings.Split(GCONFIG.MailDevs, ",")

			m = gomail.NewMessage()
			m.SetHeader("From", GCONFIG.MailEmail)
			m.SetHeader("To", arr_to...)

			if message_type == "info" {
				m.SetHeader("Subject", "hospital_ds: Информация")
			} else if message_type == "error" {
				m.SetHeader("Subject", "hospital_ds: Системная ошибка")
			} else {
				m.SetHeader("Subject", "hospital_ds: Нетипизированная ошибка")
			}

			m.SetBody("text/html", text)

			port, err = strconv.Atoi(GCONFIG.MailSmtp)

			if err != nil {
				log.Printf("panic: error = %s", err.Error())
			} else {
				d = gomail.NewDialer(GCONFIG.MailSmtp, port, GCONFIG.MailEmail, GCONFIG.MailPass)

				if err = d.DialAndSend(m); err != nil {
					log.Printf("panic: error = %s", err.Error())
				}
			}
		}
	}()
}

// Logging определение местоположения ошибки и отправки ошибки на почту
func Logging(err_text string) {
	pc, fn, line, _ := runtime.Caller(1)
	text_dev := fmt.Sprintf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err_text)

	log.Print(text_dev)
	SendDeveloper("error", text_dev)

	return
}
