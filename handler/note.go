package handler

import (
	"h24s_14-extension-backend/database"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HandleCreateNote(c echo.Context) error {
	note := new(database.Note)
	err := c.Bind(note)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to get the note")
	}

	err = database.CreateNote(note)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to insert into the database")
	}

	return c.NoContent(http.StatusOK)
}

func HandleGetNotes(c echo.Context) error {
	owner := c.Param("owner")

	notes, err := database.GetAllNotesByOwner(owner)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to get notes from the database")
	}

	return c.JSON(http.StatusOK, notes)
}

func HandleDeleteNote(c echo.Context) error {
	owner := c.QueryParam("owner")
	messageId := c.QueryParam("messageid")

	err := database.DeleteNote(owner, messageId)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to delete the note")
	}

	return c.NoContent(http.StatusOK)
}

func HandleUpdateNote(c echo.Context) error {
	note := new(database.Note)
	err := c.Bind(note)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to get the note")
	}

	err = database.UpdateNote(note.Owner, note.MessageId, note.Text, note.Color)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "failed to update the note")
	}

	return c.NoContent(http.StatusOK)
}
