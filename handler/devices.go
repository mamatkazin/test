package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"time"

	"hospital_track/env"
	"hospital_track/models"

	"github.com/gin-gonic/gin"
)

func (h *SHandler) tracks(ctx *gin.Context) {
	var (
		bodyBytes []byte
		err       error
		bodySTR   string
		body      models.SDevice
		tm        time.Time
	)

	if bodyBytes, err = ioutil.ReadAll(ctx.Request.Body); err != nil {
		env.SetAbortWithStatusJSON(err, ctx)

		return
	}

	bodySTR = string(bodyBytes)

	if err = json.Unmarshal([]byte(bodySTR), &body); err != nil {
		env.SetAbortWithStatusJSON(err, ctx)

		return
	}

	tm = time.Unix(0, body.TS*int64(time.Millisecond))
	tm = tm.Add(time.Duration(-3 * time.Hour))

	fmt.Println(tm, body.MAC, body.TS, body.Lng, body.Lat, body.Speed, body.Length*1000, math.Round(body.Hindrace), math.Round(body.Direction))

	// var input todo.TodoList
	// if err := c.BindJSON(&input); err != nil {
	// 	newErrorResponse(c, http.StatusBadRequest, err.Error())
	// 	return
	// }

	if _, err = h.services.IDevice.Registry(ctx, &body); err != nil {
		env.SetAbortWithStatusJSON(err, ctx)

		return
	}

	ctx.JSON(http.StatusOK, nil)

	return
}
