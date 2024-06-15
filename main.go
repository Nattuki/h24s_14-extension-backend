package main

import (
	"h24s_14-extension-backend/database"

	"h24s_14-extension-backend/handler"
	"h24s_14-extension-backend/middleware"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	database.UseSessionStore(e)

	withAuth := e.Group("")
	withAuth.Use(middleware.SetUserInformationToSess)

	withAuth.GET("/me", handler.HandleGetMe)

}
