package common

import (
	//	"database/sql"
	"fmt"
	//	"html"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/go-pg/pg"
	_ "github.com/nakagami/firebirdsql"
	"gopkg.in/gomail.v2"
)

type Configuration struct {
	APP_PORT        string
	MAIL_USER_EMAIL string
	MAIL_USER_PASS  string
	MAIL_SMTP       string
	MAIL_SMTP_PORT  string
	MAIL_DEVELOPERS string
	RUN_MODE        string
	PG_ADDRESS      string
	PG_USER         string
	PG_PASSWORD     string
	PG_BASE         string
}

//const G_FOLDER_PATH = "./documents/"

var (
	G_CONFIG Configuration
)

func GetPGDB(count int) (pgdb *pg.DB, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = GetRecoverError(rec)
		}
	}()

	var (
		opt pg.Options
		n   int
	)

	opt.Addr = G_CONFIG.PG_ADDRESS
	opt.Database = G_CONFIG.PG_BASE
	opt.User = G_CONFIG.PG_USER
	opt.Password = G_CONFIG.PG_PASSWORD

	if count > 10 {
		panic("Потеряно соединение с БД")
	} else {
		err = nil
	}

	pgdb = pg.Connect(&opt)

	if _, err = pgdb.QueryOne(pg.Scan(&n), "SELECT 110+1"); err != nil {
		fmt.Println("err", err)
		count = count + 1
		pgdb, err = GetPGDB(count)
	}

	return
}

func SendDeveloper(message_type, text string) {
	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				log.Print(GetRecoverErrorText(rec))

			}
		}()

		var (
			arr_to []string
			m      *gomail.Message
			d      *gomail.Dialer
			err    error
			port   int
		)

		if G_CONFIG.MAIL_DEVELOPERS != "" {
			arr_to = strings.Split(G_CONFIG.MAIL_DEVELOPERS, ",")

			m = gomail.NewMessage()
			m.SetHeader("From", G_CONFIG.MAIL_USER_EMAIL)
			m.SetHeader("To", arr_to...)

			if message_type == "info" {
				m.SetHeader("Subject", "csort_3d_ds: Информация")
			} else if message_type == "error" {
				m.SetHeader("Subject", "csort_3d_ds: Системная ошибка")
			} else {
				m.SetHeader("Subject", "csort_3d_ds: Нетипизированная ошибка")
			}

			m.SetBody("text/html", text)

			port, err = strconv.Atoi(G_CONFIG.MAIL_SMTP_PORT)

			if err != nil {
				log.Printf("panic: error = %s", err.Error())
			} else {
				d = gomail.NewDialer(G_CONFIG.MAIL_SMTP, port, G_CONFIG.MAIL_USER_EMAIL, G_CONFIG.MAIL_USER_PASS)

				if err = d.DialAndSend(m); err != nil {
					log.Printf("panic: error = %s", err.Error())
				}
			}
		}
	}()
}

func IDateToStr(date interface{}) string {
	var res string

	if date == nil {
		res = ""
	} else {
		res = date.(time.Time).Format("2006-01-02")
	}

	return res
}

func IDateTimeToStr(date interface{}) string {
	var res string

	if date == nil {
		res = ""
	} else {
		res = date.(time.Time).Format("2006-01-02 15:04:05")
	}

	return res
}

func ITimeToStr(date interface{}) string {
	var res string

	if date == nil {
		res = ""
	} else {
		res = date.(time.Time).Format("15:04")
	}

	return res
}

func ProcessingError(err_text string) (text string) {
	var (
		filename, text_dev string
		line               int
	)

	_, filename, line, _ = runtime.Caller(1)

	text = err_text
	text_dev = fmt.Sprintf("panic: %s:%d: error = %s", filepath.Base(filename), line, err_text)

	if G_CONFIG.RUN_MODE == "PRODUCT" {
		log.Print(text_dev)
		SendDeveloper("error", text_dev)
	}

	return
}

// func ProcessingInfo(info_text string) (text string) {
// 	var (
// 		filename string
// 		line     int
// 	)

// 	_, filename, line, _ = runtime.Caller(1)

// 	text = fmt.Sprintf("info: %s:%d: text = %s", filepath.Base(filename), line, info_text)

// 	log.Print(text)
// 	SendDeveloper("info", text)

// 	return
// }

func GetRecoverErrorText(rec interface{}) (text string) {
	var (
		ok  bool
		err error
	)
	err, ok = rec.(error)

	if !ok {
		err = fmt.Errorf("Непредвиденная ошибка: %v", rec)
	}

	text = err.Error()

	return
}

func GetRecoverError(rec interface{}) (err error) {
	var ok bool

	err, ok = rec.(error)

	if !ok {
		err = fmt.Errorf("Непредвиденная ошибка: %v", rec)
	}

	return
}

// //Заменяет все, что не русские или англиские буквы, и не пробел (0020) или точка (002e)
// //на знак подчеркивания (005f)
// // err always = nil
// func W1251Normalizer(inStr string) (outStr string, err error) {
// 	var validRunes = []rune{'А', 'Б', 'В', 'Г', 'Д', 'Е', 'Ж', 'З', 'И', 'Й', 'К', 'Л', 'М', 'О',
// 		'П', 'Р', 'С', 'Т', 'У', 'Ф', 'Х', 'Ц', 'Ч', 'Ш', 'Щ', 'Ъ', 'Ы', 'Ь', 'Э', 'Ю', 'Я',
// 		'а', 'б', 'в', 'г', 'д', 'е', 'ж', 'з', 'и', 'й', 'к', 'л', 'м', 'о',
// 		'п', 'р', 'с', 'т', 'у', 'ф', 'х', 'ц', 'ч', 'ш', 'щ', 'ъ', 'ы', 'ь', 'э', 'ю', 'я',
// 		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
// 		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
// 		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '(', ')',
// 		' ', '.'}
// 	var runeIsVal = func(r rune) (isVal bool) {
// 		for _, vr := range validRunes {
// 			if r == vr {
// 				isVal = true
// 				return
// 			}
// 		}
// 		isVal = false
// 		return
// 	}
// 	var inRunes = []rune(inStr)
// 	var outRunes []rune
// 	for _, rn := range inRunes {
// 		if runeIsVal(rn) {
// 			outRunes = append(outRunes, rn)
// 		} else {
// 			outRunes = append(outRunes, '_')
// 		}
// 	}
// 	return string(outRunes), nil
// }

// func CheckLenString(str string, length int) (res string) {
// 	var index int

// 	if str != "null" {
// 		str = html.EscapeString(str)

// 		if len([]rune(str)) > length {
// 			res = string([]rune(str)[0 : length-1])
// 		} else {
// 			res = str
// 		}

// 		index = strings.LastIndex(res, "&")

// 		if index > -1 {
// 			var count = len(res) - index

// 			if count < 5 {
// 				res = string([]rune(res)[0 : len(res)-count])
// 			}
// 		}

// 	} else {
// 		res = ""
// 	}

// 	return
// }
