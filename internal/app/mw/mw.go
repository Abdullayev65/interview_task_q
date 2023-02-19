package mw

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

const roleAdmin = "admin"

func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			c.String(http.StatusBadRequest, err.Error())
		}
		return nil
	}
}
