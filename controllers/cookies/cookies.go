package cookies

import (
	"net/http"

	"github.com/gorilla/securecookie"

	"nami/nami_ds/controllers/common"
)

// cookie handling

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func GetSession(request *http.Request) (sessData common.CookieData) {
	cookie, err := request.Cookie("session")

	if err == nil {
		cookieValue := make(map[string]common.CookieData)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			sessData = cookieValue["data"]
		} else {
			sessData.UserID = -1
		}
	} else {
		sessData.UserID = -1
	}

	return sessData
}

func SetSession(cookieData common.CookieData, response http.ResponseWriter) {
	value := map[string]common.CookieData{
		"data": cookieData,
	}

	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	} else {
		common.ProcessingError(err.Error())
	}
}

func ClearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)

}
