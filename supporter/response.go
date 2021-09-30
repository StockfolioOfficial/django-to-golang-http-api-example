package supporter

import "github.com/labstack/echo/v4"

func ErrorResponse(err error) echo.Map {
	return echo.Map{
		"message": err.Error(),
	}
}

func ErrorResponseMessage(msg string) echo.Map {
	return echo.Map{
		"message": msg,
	}
}