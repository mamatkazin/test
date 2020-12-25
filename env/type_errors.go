package env

import (
	"github.com/gin-gonic/gin"
)

// BisnessError структура для описания ошибки бизнес-логики
type BisnessError struct {
	Err string
}

func (e *BisnessError) Error() string {
	return e.Err
}

// HTTPError структура для описания ошибки http запросов
type HTTPError struct {
	Code int
	Err  string
}

func (e *HTTPError) Error() string {
	return e.Err
}

// SystemError структура для описания системных ошибок
type SystemError struct {
	Err string
}

func (e *SystemError) Error() string {
	return e.Err
}

func SetAbortWithStatusJSON(err error, ctx *gin.Context) {
	switch err.(type) {
	case *BisnessError:
		ctx.AbortWithStatusJSON(409, err.Error())
		return
	case *HTTPError:
		ctx.AbortWithStatusJSON(err.(*HTTPError).Code, err.(*HTTPError).Error())
		return
	case *SystemError:
		ctx.AbortWithStatusJSON(500, err.Error())
		return
	default:
		ctx.AbortWithStatusJSON(500, err.Error())
		return
	}
}
