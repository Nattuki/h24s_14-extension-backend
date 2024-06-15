package main

import (
	"h24s_14-extension-backend/database"
	"log"

	"h24s_14-extension-backend/handler"
	"h24s_14-extension-backend/middleware"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	database.UseSessionStore(e)

	withAuth := e.Group("")
	withAuth.Use(middleware.SetUserInformationToSess)
	e.GET("/loginpath", handler.HandleGetOAuthUrl)
	e.GET("/gettoken", handler.HandleGetToken)
	withAuth.GET("/me", handler.HandleGetMe)
	withAuth.GET("/note/get/:owner", handler.HandleGetNotes)
	withAuth.POST("/note/create", handler.HandleCreateNote)
	withAuth.POST("/note/update", handler.HandleUpdateNote)
	withAuth.DELETE("/note/delete", handler.HandleDeleteNote)

	err := e.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
