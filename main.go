package main

import (
	"h24s_14-extension-backend/database"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	database.UseSessionStore(e)
}
